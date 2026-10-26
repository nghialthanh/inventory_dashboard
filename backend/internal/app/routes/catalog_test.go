package routes_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"inventory_dashboard/internal/app/routes"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/database"
)

func TestCatalogListsSeededMakes(t *testing.T) {
	srv := catalogServer(t)

	rec := request(t, srv, http.MethodGet, "/api/v1/catalog/makes")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/catalog/makes status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    []struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("makes body: %v", err)
	}
	if body.Code != 200 || body.Message != "ok" {
		t.Fatalf("envelope = code %d message %q, want 200 ok", body.Code, body.Message)
	}

	got := map[string]int32{}
	var names []string
	for _, make := range body.Data {
		got[make.Name] = make.ID
		names = append(names, make.Name)
	}
	if strings.Join(names, ",") != "Ford,Honda,Toyota" {
		t.Fatalf("makes = %q, want Ford, Honda, Toyota", strings.Join(names, ","))
	}
	if got["Ford"] != 1 || got["Honda"] != 2 || got["Toyota"] != 3 {
		t.Fatalf("make addresses = %v, want Ford 1, Honda 2, Toyota 3", got)
	}
}

func TestChoosingAMakeLimitsModelsToThatMake(t *testing.T) {
	srv := catalogServer(t)
	makes := catalogData(t, srv, "/api/v1/catalog/makes")
	honda := findName(t, makes, "Honda")
	toyota := findName(t, makes, "Toyota")

	hondaModels := catalogData(t, srv, fmt.Sprintf("/api/v1/catalog/models?make_id=%d", honda.ID))
	if namesOf(hondaModels) != "Accord,Civic" {
		t.Fatalf("Honda models = %q, want Accord, Civic", namesOf(hondaModels))
	}
	for _, model := range hondaModels {
		if model.MakeID != honda.ID {
			t.Fatalf("model %s make_id = %d, want Honda %d", model.Name, model.MakeID, honda.ID)
		}
	}

	all := catalogData(t, srv, "/api/v1/catalog/models")
	if namesOf(all) != "Accord,Camry,Civic,Corolla,F-150,Mustang" {
		t.Fatalf("models = %q, want every seeded model", namesOf(all))
	}
	camry := findName(t, all, "Camry")
	if camry.MakeID != toyota.ID {
		t.Fatalf("Camry make_id = %d, want Toyota %d", camry.MakeID, toyota.ID)
	}

	unknown := catalogData(t, srv, "/api/v1/catalog/models?make_id=999999")
	if len(unknown) != 0 {
		t.Fatalf("unknown make returned %d models, want none", len(unknown))
	}

	bad := request(t, srv, http.MethodGet, "/api/v1/catalog/models?make_id=not-a-make")
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("bad make_id status = %d, want 400", bad.Code)
	}
	var errBody struct {
		Data struct {
			ErrorCode string `json:"error_code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bad.Body.Bytes(), &errBody); err != nil {
		t.Fatalf("bad make_id body: %v", err)
	}
	if errBody.Data.ErrorCode != "VALIDATION_ERROR" {
		t.Fatalf("error_code = %q, want VALIDATION_ERROR", errBody.Data.ErrorCode)
	}
	for _, raw := range []string{"0", "-1", "2147483648", "1.5"} {
		rec := request(t, srv, http.MethodGet, "/api/v1/catalog/models?make_id="+raw)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("make_id %s status = %d, want 400", raw, rec.Code)
		}
	}

	want := map[string]int32{
		"F-150": 1, "Mustang": 2, "Accord": 3, "Civic": 4, "Camry": 5, "Corolla": 6,
	}
	for _, model := range all {
		if want[model.Name] != model.ID {
			t.Fatalf("model %s address = %d, want %d", model.Name, model.ID, want[model.Name])
		}
	}
}

func TestCatalogCannotBeEdited(t *testing.T) {
	srv := catalogServer(t)
	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, path := range []string{"/api/v1/catalog/makes", "/api/v1/catalog/models"} {
			rec := request(t, srv, method, path)
			if rec.Code < 400 {
				t.Fatalf("%s %s status = %d, want a rejected catalog write", method, path, rec.Code)
			}
		}
	}
}

type catalogItem struct {
	ID     int32  `json:"id"`
	MakeID int32  `json:"make_id"`
	Name   string `json:"name"`
}

func catalogData(t *testing.T, srv handler, path string) []catalogItem {
	t.Helper()
	rec := request(t, srv, http.MethodGet, path)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d, body = %s", path, rec.Code, rec.Body.String())
	}
	var body struct {
		Code    int           `json:"code"`
		Message string        `json:"message"`
		Data    []catalogItem `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("%s body: %v", path, err)
	}
	if body.Code != 200 || body.Message != "ok" {
		t.Fatalf("%s envelope = code %d message %q", path, body.Code, body.Message)
	}
	return body.Data
}

func findName(t *testing.T, items []catalogItem, name string) catalogItem {
	t.Helper()
	for _, item := range items {
		if item.Name == name {
			return item
		}
	}
	t.Fatalf("missing %s", name)
	return catalogItem{}
}

func namesOf(items []catalogItem) string {
	names := make([]string, 0, len(items))
	for _, item := range items {
		names = append(names, item.Name)
	}
	return strings.Join(names, ",")
}

func catalogServer(t *testing.T) handler {
	t.Helper()
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := resetSchema(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE models, makes CASCADE`); err != nil {
		t.Fatalf("truncate catalog: %v", err)
	}
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("reseed: %v", err)
	}
	return routes.New(routes.Config{
		Catalog: catalog.New(repositories.NewCatalog(pool)),
	})
}

func resetSchema(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `DROP TABLE IF EXISTS vehicles, models, makes CASCADE`); err != nil {
		return err
	}
	return database.Migrate(ctx, pool)
}

func testCatalogDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()
	admin, err := pgx.Connect(ctx, "postgres://inventory:inventory@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatalf("postgres: %v", err)
	}
	defer admin.Close(ctx)

	_, err = admin.Exec(ctx, "CREATE DATABASE inventory_test")
	if err != nil && !duplicateDatabase(err) {
		t.Fatalf("create database: %v", err)
	}

	pool, err := pgxpool.New(ctx, "postgres://inventory:inventory@localhost:5432/inventory_test?sslmode=disable")
	if err != nil {
		t.Fatalf("inventory_test: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

func duplicateDatabase(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42P04"
}
