package routes_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"inventory_dashboard/internal/app/routes"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/logic/media"
	"inventory_dashboard/internal/logic/vehicle"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/database"
	"inventory_dashboard/pkg/storage"
)

func TestRecordedVehicleIsInStockWithAgeZero(t *testing.T) {
	today := time.Date(2026, 10, 26, 15, 4, 0, 0, time.UTC)
	srv := vehicleServer(t, today)

	rec := postJSON(t, srv, "/api/v1/vehicles", `{
		"vin": " 1hgcm82633a004352 ",
		"model_id": "22222222-2222-4222-8222-222222222204",
		"entered_on": "2026-10-26"
	}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", rec.Code, rec.Body.String())
	}
	created := decodeVehicle(t, rec.Body.Bytes())
	assertInStockToday(t, created)

	got := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.Data.ID).Body.Bytes())
	if request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.Data.ID).Code != http.StatusOK {
		t.Fatal("vehicle view status, want 200")
	}
	assertInStockToday(t, got)

	listRec := request(t, srv, http.MethodGet, "/api/v1/vehicles")
	if listRec.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", listRec.Code, listRec.Body.String())
	}
	var list struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Meta    struct {
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"meta"`
		Data []vehicleJSON `json:"data"`
	}
	if err := json.Unmarshal(listRec.Body.Bytes(), &list); err != nil {
		t.Fatalf("list body: %v", err)
	}
	if list.Code != 200 || list.Message != "ok" || list.Meta.Page != 1 || list.Meta.PageSize != 20 || list.Meta.Total != 1 {
		t.Fatalf("list meta = %+v message = %q, want page 1 size 20 total 1", list.Meta, list.Message)
	}
	if len(list.Data) != 1 || list.Data[0].ID != created.Data.ID {
		t.Fatalf("list = %+v, want the recorded vehicle", list.Data)
	}
	assertInStockToday(t, vehicleBody{Message: "ok", Data: list.Data[0]})
}

type vehicleJSON struct {
	ID              string  `json:"id"`
	VIN             string  `json:"vin"`
	Make            refJSON `json:"make"`
	Model           refJSON `json:"model"`
	EnteredOn       string  `json:"entered_on"`
	StockedOutOn    *string `json:"stocked_out_on"`
	StockState      string  `json:"stock_state"`
	InventoryAge    int     `json:"inventory_age"`
	IsAging         bool    `json:"is_aging"`
	ProposedAction  *string `json:"proposed_action"`
	HasPrimaryImage bool    `json:"has_primary_image"`
}

type refJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type vehicleBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    vehicleJSON `json:"data"`
}

func assertInStockToday(t *testing.T, body vehicleBody) {
	t.Helper()
	v := body.Data
	if body.Message != "ok" {
		t.Fatalf("message = %q, want ok", body.Message)
	}
	if v.VIN != "1HGCM82633A004352" {
		t.Fatalf("vin = %q, want 1HGCM82633A004352", v.VIN)
	}
	if v.Make.Name != "Honda" || v.Model.Name != "Civic" {
		t.Fatalf("make/model = %s %s, want Honda Civic", v.Make.Name, v.Model.Name)
	}
	if v.EnteredOn != "2026-10-26" || v.StockedOutOn != nil {
		t.Fatalf("entered_on = %s stocked_out_on = %v, want 2026-10-26 and empty", v.EnteredOn, v.StockedOutOn)
	}
	if v.StockState != "in_stock" || v.InventoryAge != 0 || v.IsAging {
		t.Fatalf("stock_state = %s age = %d aging = %v, want in_stock 0 false", v.StockState, v.InventoryAge, v.IsAging)
	}
	if v.ProposedAction != nil || v.HasPrimaryImage {
		t.Fatalf("proposed_action = %v has_primary_image = %v, want none and false", v.ProposedAction, v.HasPrimaryImage)
	}
}

func decodeVehicle(t *testing.T, raw []byte) vehicleBody {
	t.Helper()
	var body vehicleBody
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("vehicle body: %v", err)
	}
	return body
}

func vehicleServer(t *testing.T, today time.Time) handler {
	t.Helper()
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE vehicles`); err != nil {
		t.Fatalf("truncate vehicles: %v", err)
	}
	catalogRepo := repositories.NewCatalog(pool)
	return routes.New(routes.Config{
		Catalog:  catalog.New(catalogRepo),
		Vehicles: vehicle.New(repositories.NewVehicles(pool), catalogRepo, media.NewObjects(storage.NewMemory()), func() time.Time { return today }),
	})
}

func postJSON(t *testing.T, srv handler, path, raw string) *httptest.ResponseRecorder {
	t.Helper()
	return sendJSON(t, srv, http.MethodPost, path, raw)
}

func patchJSON(t *testing.T, srv handler, path, raw string) *httptest.ResponseRecorder {
	t.Helper()
	return sendJSON(t, srv, http.MethodPatch, path, raw)
}

func putJSON(t *testing.T, srv handler, path, raw string) *httptest.ResponseRecorder {
	t.Helper()
	return sendJSON(t, srv, http.MethodPut, path, raw)
}

func sendJSON(t *testing.T, srv handler, method, path, raw string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(raw))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}
