package routes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"inventory_dashboard/internal/app/routes"
	"inventory_dashboard/internal/logic/catalog"
	"inventory_dashboard/internal/logic/vehicle"
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/pkg/database"
)

const (
	fordID    = "11111111-1111-4111-8111-111111111101"
	hondaID   = "11111111-1111-4111-8111-111111111102"
	f150ID    = "22222222-2222-4222-8222-222222222201"
	civicID   = "22222222-2222-4222-8222-222222222204"
	civicVIN  = "1HGCM82633A004352"
	secondVIN = "2HGFG3B54GH123456"
	thirdVIN  = "3N1AB7AP5HY123456"
)

func TestInventoryAgeBoundariesAndFutureEnteredOn(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)

	ninety := recordVehicle(t, srv, civicVIN, civicID, "2026-07-28")
	if ninety.InventoryAge != 90 || ninety.IsAging || ninety.StockState != "in_stock" {
		t.Fatalf("90 days = age %d aging %v state %s, want 90 false in_stock", ninety.InventoryAge, ninety.IsAging, ninety.StockState)
	}

	ninetyOne := recordVehicle(t, srv, secondVIN, civicID, "2026-07-27")
	if ninetyOne.InventoryAge != 91 || !ninetyOne.IsAging {
		t.Fatalf("91 days = age %d aging %v, want 91 true", ninetyOne.InventoryAge, ninetyOne.IsAging)
	}

	future := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(thirdVIN, civicID, "2026-10-27"))
	assertAPIError(t, future, http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestRejectedVINAndMissingFields(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))

	cases := []string{
		`{"vin":"1HGCM82633A00435I","model_id":"` + civicID + `","entered_on":"2026-10-26"}`,
		`{"vin":"1HGCM82633A00435O","model_id":"` + civicID + `","entered_on":"2026-10-26"}`,
		`{"vin":"1HGCM82633A00435Q","model_id":"` + civicID + `","entered_on":"2026-10-26"}`,
		`{"vin":"SHORT","model_id":"` + civicID + `","entered_on":"2026-10-26"}`,
		`{"model_id":"` + civicID + `","entered_on":"2026-10-26"}`,
		`{"vin":"` + civicVIN + `","entered_on":"2026-10-26"}`,
		`{"vin":"` + civicVIN + `","model_id":"` + civicID + `"}`,
		`{"vin":"` + civicVIN + `","model_id":"33333333-3333-4333-8333-333333333333","entered_on":"2026-10-26"}`,
		`{"vin":"not-a-uuid","model_id":"not-a-uuid","entered_on":"yesterday"}`,
	}
	for _, body := range cases {
		rec := postJSON(t, srv, "/api/v1/vehicles", body)
		assertAPIError(t, rec, http.StatusBadRequest, "VALIDATION_ERROR")
	}
}

func TestDuplicateVINTellsTheManagerItIsAlreadyUsed(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))
	recordVehicle(t, srv, civicVIN, civicID, "2026-10-01")

	again := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(civicVIN, f150ID, "2026-10-02"))
	assertAPIError(t, again, http.StatusConflict, "CONFLICT")
	var body errorJSON
	if err := json.Unmarshal(again.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Message != "VIN is already used" {
		t.Fatalf("message = %q, want VIN is already used", body.Message)
	}

	list := listVehicles(t, srv, "")
	if list.Meta.Total != 1 {
		t.Fatalf("total = %d, want 1", list.Meta.Total)
	}
}

func TestSimultaneousRecordsOfTheSameVINCreateOneVehicle(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))
	raw := vehicleBodyJSON("  "+civicVIN+" ", civicID, "2026-10-26")

	var wg sync.WaitGroup
	recs := make([]*httptest.ResponseRecorder, 2)
	wg.Add(2)
	for i := range recs {
		go func(i int) {
			defer wg.Done()
			recs[i] = postJSON(t, srv, "/api/v1/vehicles", raw)
		}(i)
	}
	wg.Wait()

	created, conflicts := 0, 0
	for _, rec := range recs {
		switch rec.Code {
		case http.StatusCreated:
			created++
		case http.StatusConflict:
			conflicts++
			assertAPIError(t, rec, http.StatusConflict, "CONFLICT")
		default:
			t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
		}
	}
	if created != 1 || conflicts != 1 {
		t.Fatalf("created = %d conflicts = %d, want 1 and 1", created, conflicts)
	}
	if listVehicles(t, srv, "").Meta.Total != 1 {
		t.Fatal("simultaneous duplicate VIN created more than one vehicle")
	}
}

func TestStockedOutStayStopsAgeAndKeepsTheVIN(t *testing.T) {
	today := time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC)
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE vehicles`); err != nil {
		t.Fatalf("truncate vehicles: %v", err)
	}
	const id = "55555555-5555-4555-8555-555555555555"
	if _, err := pool.Exec(ctx, `
		INSERT INTO vehicles (id, vin, model_id, entered_on, stocked_out_on)
		VALUES ($1, $2, $3, DATE '2026-07-01', DATE '2026-07-11')
	`, id, civicVIN, civicID); err != nil {
		t.Fatalf("insert stocked-out vehicle: %v", err)
	}

	catalogRepo := repositories.NewCatalog(pool)
	srv := routes.New(routes.Config{
		Catalog:  catalog.New(catalogRepo),
		Vehicles: vehicle.New(repositories.NewVehicles(pool), catalogRepo, func() time.Time { return today }),
	})
	rec := request(t, srv, http.MethodGet, "/api/v1/vehicles/"+id)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	got := decodeVehicle(t, rec.Body.Bytes())
	if got.Data.StockState != "stocked_out" || got.Data.InventoryAge != 10 || got.Data.IsAging {
		t.Fatalf("stocked-out stay = state %s age %d aging %v, want stocked_out 10 false", got.Data.StockState, got.Data.InventoryAge, got.Data.IsAging)
	}
	if got.Data.StockedOutOn == nil || *got.Data.StockedOutOn != "2026-07-11" {
		t.Fatalf("stocked_out_on = %v, want 2026-07-11", got.Data.StockedOutOn)
	}

	again := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(civicVIN, f150ID, "2026-10-26"))
	assertAPIError(t, again, http.StatusConflict, "CONFLICT")
}

func TestUnknownVehicleAddress(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))
	unknown := request(t, srv, http.MethodGet, "/api/v1/vehicles/44444444-4444-4444-8444-444444444444")
	assertAPIError(t, unknown, http.StatusNotFound, "NOT_FOUND")

	bad := request(t, srv, http.MethodGet, "/api/v1/vehicles/not-a-vehicle")
	assertAPIError(t, bad, http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestListOrderPagingAndEmptyInventory(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))

	empty := listVehicles(t, srv, "")
	if empty.Meta.Total != 0 || empty.Meta.Page != 1 || empty.Meta.PageSize != 20 || len(empty.Data) != 0 {
		t.Fatalf("empty list = %+v", empty)
	}
	if empty.RawData != "[]" {
		t.Fatalf("empty data = %s, want []", empty.RawData)
	}

	later := recordVehicle(t, srv, civicVIN, civicID, "2026-10-20")
	earlyA := recordVehicle(t, srv, secondVIN, f150ID, "2026-10-01")
	earlyB := recordVehicle(t, srv, thirdVIN, civicID, "2026-10-01")

	first := earlyA.ID
	second := earlyB.ID
	if first > second {
		first, second = second, first
	}
	list := listVehicles(t, srv, "")
	if len(list.Data) != 3 || list.Data[0].ID != first || list.Data[1].ID != second || list.Data[2].ID != later.ID {
		t.Fatalf("order = %s %s %s, want %s %s %s", list.Data[0].ID, list.Data[1].ID, list.Data[2].ID, first, second, later.ID)
	}

	page := listVehicles(t, srv, "?page_size=2")
	if page.Meta.PageSize != 2 || page.Meta.Total != 3 || len(page.Data) != 2 || page.Data[0].ID != first {
		t.Fatalf("first page = %+v", page)
	}
	next := listVehicles(t, srv, "?page=2&page_size=2")
	if len(next.Data) != 1 || next.Data[0].ID != later.ID || next.Meta.Total != 3 {
		t.Fatalf("second page = %+v", next)
	}
	past := listVehicles(t, srv, "?page=3&page_size=2")
	if past.Meta.Total != 3 || past.RawData != "[]" {
		t.Fatalf("page past the end = total %d data %s", past.Meta.Total, past.RawData)
	}

	capped := listVehicles(t, srv, "?page_size=250")
	if capped.Meta.PageSize != 100 || capped.Meta.Total != 3 {
		t.Fatalf("capped page size = %d total %d, want 100 and 3", capped.Meta.PageSize, capped.Meta.Total)
	}
	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?page=0"), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?page_size=0"), http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestAgingOnlyFilterAndInventoryAgeBoundaries(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE vehicles`); err != nil {
		t.Fatalf("truncate vehicles: %v", err)
	}
	const stayedID = "55555555-5555-4555-8555-555555555505"
	const stayedVIN = "5NPE24AF5FH123456"
	if _, err := pool.Exec(ctx, `
		INSERT INTO vehicles (id, vin, model_id, entered_on, stocked_out_on)
		VALUES ($1, $2, $3, DATE '2026-07-01', DATE '2026-10-10')
	`, stayedID, stayedVIN, civicID); err != nil {
		t.Fatalf("insert stocked-out vehicle: %v", err)
	}
	catalogRepo := repositories.NewCatalog(pool)
	srv := routes.New(routes.Config{
		Catalog:  catalog.New(catalogRepo),
		Vehicles: vehicle.New(repositories.NewVehicles(pool), catalogRepo, func() time.Time { return today }),
	})

	age0 := recordVehicle(t, srv, civicVIN, civicID, "2026-10-26")
	age90 := recordVehicle(t, srv, secondVIN, f150ID, "2026-07-28")
	age91 := recordVehicle(t, srv, thirdVIN, civicID, "2026-07-27")

	all := listVehicles(t, srv, "")
	if all.Meta.Total != 4 || len(all.Data) != 4 {
		t.Fatalf("unfiltered total = %d len = %d, want 4 and 4", all.Meta.Total, len(all.Data))
	}
	if all.Data[0].ID != stayedID || all.Data[0].StockState != "stocked_out" || all.Data[0].InventoryAge != 101 || all.Data[0].IsAging {
		t.Fatalf("stocked-out stay = %+v, want stocked_out age 101 not aging", all.Data[0])
	}
	if all.Data[1].ID != age91.ID || all.Data[1].StockState != "in_stock" || all.Data[1].InventoryAge != 91 || !all.Data[1].IsAging {
		t.Fatalf("91 days = %+v, want in_stock aging stock", all.Data[1])
	}
	if all.Data[2].ID != age90.ID || all.Data[2].StockState != "in_stock" || all.Data[2].InventoryAge != 90 || all.Data[2].IsAging {
		t.Fatalf("90 days = %+v, want in_stock and not aging", all.Data[2])
	}
	if all.Data[3].ID != age0.ID || all.Data[3].StockState != "in_stock" || all.Data[3].InventoryAge != 0 || all.Data[3].IsAging {
		t.Fatalf("0 days = %+v, want in_stock and not aging", all.Data[3])
	}

	only := listVehicles(t, srv, "?aging=true")
	if only.Meta.Total != 1 || len(only.Data) != 1 || only.Data[0].ID != age91.ID || only.Data[0].StockState != "in_stock" || !only.Data[0].IsAging {
		t.Fatalf("aging only = total %d data %+v, want the 91-day vehicle", only.Meta.Total, only.Data)
	}

	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?aging=false"), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?aging=yes"), http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestCorrectingAnInStockVehicleKeepsItsAddress(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-10-26")

	rec := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{
		"vin": " 2hgfg3b54gh123456 ",
		"model_id": "`+f150ID+`",
		"entered_on": "2026-07-27"
	}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("correct status = %d, body = %s", rec.Code, rec.Body.String())
	}
	corrected := decodeVehicle(t, rec.Body.Bytes())
	assertCorrected(t, corrected, created.ID)

	view := request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID)
	if view.Code != http.StatusOK {
		t.Fatalf("view status = %d, body = %s", view.Code, view.Body.String())
	}
	assertCorrected(t, decodeVehicle(t, view.Body.Bytes()), created.ID)
}

func assertCorrected(t *testing.T, body vehicleBody, id string) {
	t.Helper()
	v := body.Data
	if body.Message != "ok" || v.ID != id {
		t.Fatalf("address = %s message = %q, want %s and ok", v.ID, body.Message, id)
	}
	if v.VIN != secondVIN || v.Make.Name != "Ford" || v.Model.Name != "F-150" {
		t.Fatalf("identity = %s %s %s, want %s Ford F-150", v.VIN, v.Make.Name, v.Model.Name, secondVIN)
	}
	if v.EnteredOn != "2026-07-27" || v.StockedOutOn != nil || v.StockState != "in_stock" || v.InventoryAge != 91 || !v.IsAging {
		t.Fatalf("stay = entered %s state %s age %d aging %v, want 2026-07-27 in_stock 91 true", v.EnteredOn, v.StockState, v.InventoryAge, v.IsAging)
	}
	if v.ProposedAction != nil {
		t.Fatalf("proposed_action = %v, want unchanged empty value", v.ProposedAction)
	}
}

func TestCorrectionKeepsASubsetAndRejectsOtherCommands(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-10-01")

	onlyVIN := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"vin":" `+secondVIN+` "}`)
	if onlyVIN.Code != http.StatusOK {
		t.Fatalf("vin correction status = %d, body = %s", onlyVIN.Code, onlyVIN.Body.String())
	}
	vinOnly := decodeVehicle(t, onlyVIN.Body.Bytes()).Data
	if vinOnly.ID != created.ID || vinOnly.VIN != secondVIN || vinOnly.Model.Name != "Civic" || vinOnly.EnteredOn != "2026-10-01" {
		t.Fatalf("vin-only correction = %+v", vinOnly)
	}

	onlyModel := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"model_id":"`+f150ID+`"}`)
	if onlyModel.Code != http.StatusOK {
		t.Fatalf("model correction status = %d, body = %s", onlyModel.Code, onlyModel.Body.String())
	}
	modelOnly := decodeVehicle(t, onlyModel.Body.Bytes()).Data
	if modelOnly.ID != created.ID || modelOnly.VIN != secondVIN || modelOnly.Make.Name != "Ford" || modelOnly.Model.Name != "F-150" || modelOnly.EnteredOn != "2026-10-01" {
		t.Fatalf("model-only correction = %+v", modelOnly)
	}

	onlyEntered := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"entered_on":"2026-07-28"}`)
	if onlyEntered.Code != http.StatusOK {
		t.Fatalf("entered-on correction status = %d, body = %s", onlyEntered.Code, onlyEntered.Body.String())
	}
	enteredOnly := decodeVehicle(t, onlyEntered.Body.Bytes()).Data
	if enteredOnly.ID != created.ID || enteredOnly.VIN != secondVIN || enteredOnly.Model.Name != "F-150" || enteredOnly.EnteredOn != "2026-07-28" || enteredOnly.InventoryAge != 90 || enteredOnly.IsAging || enteredOnly.StockState != "in_stock" {
		t.Fatalf("entered-on-only correction = %+v", enteredOnly)
	}

	mixed := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{
		"entered_on": "2026-07-28",
		"proposed_action": "hold",
		"stocked_out_on": "2026-10-01"
	}`)
	assertAPIError(t, mixed, http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{}`), http.StatusBadRequest, "VALIDATION_ERROR")

	kept := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if kept.ID != created.ID || kept.VIN != secondVIN || kept.EnteredOn != "2026-07-28" || kept.Model.Name != "F-150" {
		t.Fatalf("vehicle changed after a rejected correction: %+v", kept)
	}
	if kept.StockedOutOn != nil || kept.ProposedAction != nil || kept.StockState != "in_stock" {
		t.Fatalf("forbidden fields changed: %+v", kept)
	}
}

func TestCorrectedVINStaysUniqueAndEnteredOnCannotBeAfterToday(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE vehicles`); err != nil {
		t.Fatalf("truncate vehicles: %v", err)
	}
	const stayedID = "55555555-5555-4555-8555-555555555506"
	if _, err := pool.Exec(ctx, `
		INSERT INTO vehicles (id, vin, model_id, entered_on, stocked_out_on)
		VALUES ($1, $2, $3, DATE '2026-07-01', DATE '2026-07-11')
	`, stayedID, thirdVIN, civicID); err != nil {
		t.Fatalf("insert stocked-out vehicle: %v", err)
	}
	catalogRepo := repositories.NewCatalog(pool)
	srv := routes.New(routes.Config{
		Catalog:  catalog.New(catalogRepo),
		Vehicles: vehicle.New(repositories.NewVehicles(pool), catalogRepo, func() time.Time { return today }),
	})
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-10-26")
	recordVehicle(t, srv, secondVIN, f150ID, "2026-10-02")

	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"entered_on":"2026-10-27"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"vin":"1HGCM82633A00435I"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"model_id":"33333333-3333-4333-8333-333333333333"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"vin":"`+secondVIN+`"}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"vin":"`+thirdVIN+`"}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, patchJSON(t, srv, "/api/v1/vehicles/"+stayedID, `{"vin":"JH4KA4650MC123456"}`), http.StatusConflict, "CONFLICT")

	kept := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if kept.VIN != civicVIN || kept.EnteredOn != "2026-10-26" || kept.Model.ID != civicID || kept.StockState != "in_stock" {
		t.Fatalf("in-stock vehicle changed after rejected corrections: %+v", kept)
	}
	stayed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+stayedID).Body.Bytes()).Data
	if stayed.VIN != thirdVIN || stayed.StockState != "stocked_out" || stayed.InventoryAge != 10 || stayed.IsAging {
		t.Fatalf("stocked-out vehicle changed: %+v", stayed)
	}
}

func TestMakeAndModelFilters(t *testing.T) {
	srv := vehicleServer(t, time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC))
	ford := recordVehicle(t, srv, civicVIN, f150ID, "2026-10-01")
	honda := recordVehicle(t, srv, secondVIN, civicID, "2026-10-02")

	byMake := listVehicles(t, srv, "?make_id="+fordID)
	if byMake.Meta.Total != 1 || byMake.Data[0].ID != ford.ID || byMake.Data[0].Make.Name != "Ford" {
		t.Fatalf("make filter = %+v", byMake.Data)
	}
	byModel := listVehicles(t, srv, "?model_id="+civicID)
	if byModel.Meta.Total != 1 || byModel.Data[0].ID != honda.ID || byModel.Data[0].Model.Name != "Civic" {
		t.Fatalf("model filter = %+v", byModel.Data)
	}
	both := listVehicles(t, srv, "?make_id="+hondaID+"&model_id="+civicID)
	if both.Meta.Total != 1 || both.Data[0].ID != honda.ID {
		t.Fatalf("make and model filter = %+v", both.Data)
	}

	mismatch := request(t, srv, http.MethodGet, "/api/v1/vehicles?make_id="+fordID+"&model_id="+civicID)
	assertAPIError(t, mismatch, http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?make_id=bad"), http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestStockOutKeepsTheVehicleListedAsStockedOut(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-07-01")

	rec := postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-10-10"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("stock-out status = %d, body = %s", rec.Code, rec.Body.String())
	}
	departed := decodeVehicle(t, rec.Body.Bytes())
	if departed.Message != "ok" {
		t.Fatalf("message = %q, want ok", departed.Message)
	}
	v := departed.Data
	if v.ID != created.ID || v.VIN != civicVIN {
		t.Fatalf("identity = %s %s, want %s %s", v.ID, v.VIN, created.ID, civicVIN)
	}
	if v.StockedOutOn == nil || *v.StockedOutOn != "2026-10-10" || v.StockState != "stocked_out" {
		t.Fatalf("departure = stocked_out_on %v state %s, want 2026-10-10 stocked_out", v.StockedOutOn, v.StockState)
	}
	if v.InventoryAge != 101 || v.IsAging {
		t.Fatalf("stay = age %d aging %v, want 101 false", v.InventoryAge, v.IsAging)
	}

	list := listVehicles(t, srv, "")
	if list.Meta.Total != 1 || len(list.Data) != 1 || list.Data[0].ID != created.ID {
		t.Fatalf("list after stock-out = %+v, want the retained vehicle", list)
	}
	listed := list.Data[0]
	if listed.StockState != "stocked_out" || listed.InventoryAge != 101 || listed.IsAging {
		t.Fatalf("listed stay = %+v, want stocked_out age 101 not aging", listed)
	}
	if listed.StockedOutOn == nil || *listed.StockedOutOn != "2026-10-10" {
		t.Fatalf("listed stocked_out_on = %v, want 2026-10-10", listed.StockedOutOn)
	}
}

func TestStockOutRejectsBadDatesAndUnknownVehicles(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-10-01")

	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-09-30"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-10-27"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":""}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/44444444-4444-4444-8444-444444444444/stock-out", `{"stocked_out_on":"2026-10-26"}`), http.StatusNotFound, "NOT_FOUND")

	kept := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if kept.StockState != "in_stock" || kept.StockedOutOn != nil {
		t.Fatalf("vehicle changed after rejected stock-out: %+v", kept)
	}
}

func TestStockOutVINStaysUsedAndDepartureCannotReopen(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-07-01")

	rec := postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-10-10"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("stock-out status = %d, body = %s", rec.Code, rec.Body.String())
	}

	again := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(civicVIN, f150ID, "2026-10-26"))
	assertAPIError(t, again, http.StatusConflict, "CONFLICT")
	var conflict errorJSON
	if err := json.Unmarshal(again.Body.Bytes(), &conflict); err != nil {
		t.Fatal(err)
	}
	if conflict.Message != "VIN is already used" {
		t.Fatalf("message = %q, want VIN is already used", conflict.Message)
	}

	reopen := postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":null}`)
	assertAPIError(t, reopen, http.StatusConflict, "CONFLICT")
	var body errorJSON
	if err := json.Unmarshal(reopen.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Message != "The departure stands." {
		t.Fatalf("message = %q, want The departure stands.", body.Message)
	}

	stayed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if stayed.StockState != "stocked_out" || stayed.StockedOutOn == nil || *stayed.StockedOutOn != "2026-10-10" || stayed.InventoryAge != 101 || stayed.IsAging {
		t.Fatalf("after reopen attempt = %+v, want stocked_out 2026-10-10 age 101 not aging", stayed)
	}
}

func TestDemoStockedOutVehicleIsPresentOnFirstOpen(t *testing.T) {
	today := time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)
	ctx := context.Background()
	pool := testCatalogDB(t)
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE vehicles`); err != nil {
		t.Fatalf("truncate vehicles: %v", err)
	}
	if err := database.Migrate(ctx, pool); err != nil {
		t.Fatalf("reseed: %v", err)
	}
	const demoID = "55555555-5555-4555-8555-555555555504"
	// Seed dates follow CURRENT_TIMESTAMP; pin the demo stay to the fixed clock.
	if _, err := pool.Exec(ctx, `
		UPDATE vehicles
		SET entered_on = DATE '2026-06-28',
			stocked_out_on = DATE '2026-10-06'
		WHERE id = $1
	`, demoID); err != nil {
		t.Fatalf("align demo stay: %v", err)
	}
	catalogRepo := repositories.NewCatalog(pool)
	srv := routes.New(routes.Config{
		Catalog:  catalog.New(catalogRepo),
		Vehicles: vehicle.New(repositories.NewVehicles(pool), catalogRepo, func() time.Time { return today }),
	})

	rec := request(t, srv, http.MethodGet, "/api/v1/vehicles/"+demoID)
	if rec.Code != http.StatusOK {
		t.Fatalf("demo status = %d, body = %s", rec.Code, rec.Body.String())
	}
	demo := decodeVehicle(t, rec.Body.Bytes()).Data
	if demo.StockState != "stocked_out" || demo.StockedOutOn == nil || *demo.StockedOutOn != "2026-10-06" || demo.IsAging {
		t.Fatalf("demo = %+v, want stocked_out 2026-10-06 not aging", demo)
	}
	if demo.InventoryAge != 100 {
		t.Fatalf("demo inventory_age = %d, want 100 (stay longer than 90 days)", demo.InventoryAge)
	}

	list := listVehicles(t, srv, "")
	found := false
	for _, item := range list.Data {
		if item.ID == demoID {
			found = true
			if item.StockState != "stocked_out" || item.IsAging || item.InventoryAge != 100 {
				t.Fatalf("listed demo = %+v, want stocked_out age 100 not aging", item)
			}
		}
	}
	if !found {
		t.Fatal("demo stocked-out vehicle missing from the list on first open")
	}
}

func TestListFiltersByStockStateAndStockedOutOn(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)

	inStock := recordVehicle(t, srv, civicVIN, civicID, "2026-10-01")
	aging := recordVehicle(t, srv, secondVIN, f150ID, "2026-07-27")
	dayA := recordVehicle(t, srv, thirdVIN, civicID, "2026-07-01")
	stockOutVehicle(t, srv, dayA.ID, "2026-10-10")
	dayB := recordVehicle(t, srv, "JH4KA4650MC123456", civicID, "2026-08-01")
	stockOutVehicle(t, srv, dayB.ID, "2026-10-15")

	onlyIn := listVehicles(t, srv, "?stock_state=in_stock")
	if onlyIn.Meta.Total != 2 || len(onlyIn.Data) != 2 {
		t.Fatalf("in_stock filter = total %d data %+v, want 2", onlyIn.Meta.Total, onlyIn.Data)
	}
	for _, item := range onlyIn.Data {
		if item.StockState != "in_stock" {
			t.Fatalf("in_stock filter included %+v", item)
		}
	}

	onlyOut := listVehicles(t, srv, "?stock_state=stocked_out")
	if onlyOut.Meta.Total != 2 || len(onlyOut.Data) != 2 {
		t.Fatalf("stocked_out filter = total %d data %+v, want 2", onlyOut.Meta.Total, onlyOut.Data)
	}
	for _, item := range onlyOut.Data {
		if item.StockState != "stocked_out" {
			t.Fatalf("stocked_out filter included %+v", item)
		}
	}

	onDay := listVehicles(t, srv, "?stocked_out_on=2026-10-10")
	if onDay.Meta.Total != 1 || len(onDay.Data) != 1 || onDay.Data[0].ID != dayA.ID {
		t.Fatalf("stocked_out_on filter = %+v, want vehicle %s", onDay.Data, dayA.ID)
	}
	if onDay.Data[0].StockedOutOn == nil || *onDay.Data[0].StockedOutOn != "2026-10-10" {
		t.Fatalf("stocked_out_on = %v, want 2026-10-10", onDay.Data[0].StockedOutOn)
	}

	hondaOut := listVehicles(t, srv, "?make_id="+hondaID+"&stock_state=stocked_out")
	if hondaOut.Meta.Total != 2 || len(hondaOut.Data) != 2 {
		t.Fatalf("make + stocked_out = %+v, want both Honda departures", hondaOut.Data)
	}

	civicDay := listVehicles(t, srv, "?model_id="+civicID+"&stocked_out_on=2026-10-15")
	if civicDay.Meta.Total != 1 || civicDay.Data[0].ID != dayB.ID {
		t.Fatalf("model + stocked_out_on = %+v, want %s", civicDay.Data, dayB.ID)
	}

	agingOnly := listVehicles(t, srv, "?aging=true&stock_state=in_stock")
	if agingOnly.Meta.Total != 1 || agingOnly.Data[0].ID != aging.ID || !agingOnly.Data[0].IsAging {
		t.Fatalf("aging + in_stock = %+v, want aging vehicle %s", agingOnly.Data, aging.ID)
	}
	agingOut := listVehicles(t, srv, "?aging=true&stock_state=stocked_out")
	if agingOut.Meta.Total != 0 || len(agingOut.Data) != 0 {
		t.Fatalf("aging + stocked_out = %+v, want empty", agingOut.Data)
	}

	kept := listVehicles(t, srv, "")
	if kept.Meta.Total != 4 {
		t.Fatalf("unfiltered total = %d, want 4 (in-stock %s still present)", kept.Meta.Total, inStock.ID)
	}

	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?stock_state=sold"), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, request(t, srv, http.MethodGet, "/api/v1/vehicles?stocked_out_on=10-10-2026"), http.StatusBadRequest, "VALIDATION_ERROR")
}

func TestCorrectionAfterStockOutIsRejectedAndFieldsStay(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-07-01")
	stockOutVehicle(t, srv, created.ID, "2026-10-10")

	attempts := []string{
		`{"vin":"` + secondVIN + `"}`,
		`{"model_id":"` + f150ID + `"}`,
		`{"entered_on":"2026-08-01"}`,
		`{"vin":"` + secondVIN + `","model_id":"` + f150ID + `","entered_on":"2026-08-01"}`,
	}
	for _, body := range attempts {
		rec := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, body)
		assertAPIError(t, rec, http.StatusConflict, "CONFLICT")
		var errBody errorJSON
		if err := json.Unmarshal(rec.Body.Bytes(), &errBody); err != nil {
			t.Fatal(err)
		}
		if errBody.Message != "The correction is rejected because the vehicle is stocked out." {
			t.Fatalf("message = %q for body %s", errBody.Message, body)
		}
	}

	stayed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if stayed.VIN != civicVIN || stayed.Model.ID != civicID || stayed.EnteredOn != "2026-07-01" {
		t.Fatalf("identity changed after rejected corrections: %+v", stayed)
	}
	if stayed.StockState != "stocked_out" || stayed.StockedOutOn == nil || *stayed.StockedOutOn != "2026-10-10" {
		t.Fatalf("stock state changed after rejected corrections: %+v", stayed)
	}
	if stayed.InventoryAge != 101 || stayed.IsAging {
		t.Fatalf("stay = age %d aging %v, want 101 false", stayed.InventoryAge, stayed.IsAging)
	}
}

func TestProposedActionOnAgingStockIsSetAndShown(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	aging := recordVehicle(t, srv, civicVIN, civicID, "2026-07-27")
	if aging.InventoryAge != 91 || !aging.IsAging {
		t.Fatalf("setup = age %d aging %v, want 91 true", aging.InventoryAge, aging.IsAging)
	}

	rec := putJSON(t, srv, "/api/v1/vehicles/"+aging.ID+"/proposed-action", `{"proposed_action":"price_reduction_planned"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("set status = %d, body = %s", rec.Code, rec.Body.String())
	}
	set := decodeVehicle(t, rec.Body.Bytes()).Data
	if set.ProposedAction == nil || *set.ProposedAction != "price_reduction_planned" {
		t.Fatalf("proposed_action = %v, want price_reduction_planned", set.ProposedAction)
	}
	if !set.IsAging || set.StockState != "in_stock" {
		t.Fatalf("after set = aging %v state %s, want aging in_stock", set.IsAging, set.StockState)
	}

	viewed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+aging.ID).Body.Bytes()).Data
	if viewed.ProposedAction == nil || *viewed.ProposedAction != "price_reduction_planned" {
		t.Fatalf("view proposed_action = %v, want price_reduction_planned", viewed.ProposedAction)
	}
}

func TestHoldIsAProposedActionDistinctFromNone(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	aging := recordVehicle(t, srv, civicVIN, civicID, "2026-07-27")
	if aging.ProposedAction != nil {
		t.Fatalf("start proposed_action = %v, want none", aging.ProposedAction)
	}

	held := setProposedAction(t, srv, aging.ID, "hold")
	if held.ProposedAction == nil || *held.ProposedAction != "hold" {
		t.Fatalf("hold = %v, want hold (not none)", held.ProposedAction)
	}
}

func TestProposedActionReplacesAndEmptyClears(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	aging := recordVehicle(t, srv, civicVIN, civicID, "2026-07-27")

	first := setProposedAction(t, srv, aging.ID, "price_reduction_planned")
	if first.ProposedAction == nil || *first.ProposedAction != "price_reduction_planned" {
		t.Fatalf("first = %v, want price_reduction_planned", first.ProposedAction)
	}

	replaced := setProposedAction(t, srv, aging.ID, "wholesale_planned")
	if replaced.ProposedAction == nil || *replaced.ProposedAction != "wholesale_planned" {
		t.Fatalf("replaced = %v, want wholesale_planned", replaced.ProposedAction)
	}

	cleared := setProposedAction(t, srv, aging.ID, "")
	if cleared.ProposedAction != nil {
		t.Fatalf("cleared = %v, want none", cleared.ProposedAction)
	}
}

func TestProposedActionRejectedWhenNotAgingStock(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)

	fresh := recordVehicle(t, srv, civicVIN, civicID, "2026-10-26")
	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+fresh.ID+"/proposed-action", `{"proposed_action":"hold"}`), http.StatusConflict, "CONFLICT")

	ninety := recordVehicle(t, srv, secondVIN, civicID, "2026-07-28")
	if ninety.InventoryAge != 90 || ninety.IsAging {
		t.Fatalf("90 days = age %d aging %v, want 90 false", ninety.InventoryAge, ninety.IsAging)
	}
	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+ninety.ID+"/proposed-action", `{"proposed_action":"hold"}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+ninety.ID+"/proposed-action", `{"proposed_action":""}`), http.StatusConflict, "CONFLICT")

	kept := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+ninety.ID).Body.Bytes()).Data
	if kept.ProposedAction != nil {
		t.Fatalf("proposed_action after reject = %v, want none", kept.ProposedAction)
	}
}

func TestProposedActionRejectsUnknownIntentions(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	aging := recordVehicle(t, srv, civicVIN, civicID, "2026-07-27")

	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+aging.ID+"/proposed-action", `{"proposed_action":"discount"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+aging.ID+"/proposed-action", `{"proposed_action":"no_action"}`), http.StatusBadRequest, "VALIDATION_ERROR")

	kept := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+aging.ID).Body.Bytes()).Data
	if kept.ProposedAction != nil {
		t.Fatalf("proposed_action after reject = %v, want none", kept.ProposedAction)
	}
}

func TestProposedActionStaysVisibleAfterStockOutAndFurtherChangeIsRejected(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	aging := recordVehicle(t, srv, civicVIN, civicID, "2026-07-27")
	setProposedAction(t, srv, aging.ID, "wholesale_planned")

	departed := stockOutVehicle(t, srv, aging.ID, "2026-10-20")
	if departed.ProposedAction == nil || *departed.ProposedAction != "wholesale_planned" {
		t.Fatalf("after stock-out proposed_action = %v, want wholesale_planned", departed.ProposedAction)
	}
	if departed.IsAging || departed.StockState != "stocked_out" {
		t.Fatalf("after stock-out = aging %v state %s, want stocked_out not aging", departed.IsAging, departed.StockState)
	}

	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+aging.ID+"/proposed-action", `{"proposed_action":"hold"}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, putJSON(t, srv, "/api/v1/vehicles/"+aging.ID+"/proposed-action", `{"proposed_action":""}`), http.StatusConflict, "CONFLICT")

	stayed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+aging.ID).Body.Bytes()).Data
	if stayed.ProposedAction == nil || *stayed.ProposedAction != "wholesale_planned" {
		t.Fatalf("after rejected change = %v, want wholesale_planned", stayed.ProposedAction)
	}
}

func setProposedAction(t *testing.T, srv handler, id, value string) vehicleJSON {
	t.Helper()
	raw, err := json.Marshal(map[string]string{"proposed_action": value})
	if err != nil {
		t.Fatal(err)
	}
	rec := putJSON(t, srv, "/api/v1/vehicles/"+id+"/proposed-action", string(raw))
	if rec.Code != http.StatusOK {
		t.Fatalf("proposed-action status = %d, body = %s", rec.Code, rec.Body.String())
	}
	return decodeVehicle(t, rec.Body.Bytes()).Data
}

func TestLaterStockOutReplacesTheDateAndClearingIsRejected(t *testing.T) {
	today := time.Date(2026, 10, 26, 8, 0, 0, 0, time.UTC)
	srv := vehicleServer(t, today)
	created := recordVehicle(t, srv, civicVIN, civicID, "2026-07-01")
	stockOutVehicle(t, srv, created.ID, "2026-10-10")

	replaced := stockOutVehicle(t, srv, created.ID, "2026-10-20")
	if replaced.StockedOutOn == nil || *replaced.StockedOutOn != "2026-10-20" || replaced.StockState != "stocked_out" {
		t.Fatalf("replaced = %+v, want stocked_out_on 2026-10-20", replaced)
	}
	if replaced.InventoryAge != 111 || replaced.IsAging {
		t.Fatalf("replaced stay = age %d aging %v, want 111 false", replaced.InventoryAge, replaced.IsAging)
	}

	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":null}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":""}`), http.StatusConflict, "CONFLICT")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-06-30"}`), http.StatusBadRequest, "VALIDATION_ERROR")
	assertAPIError(t, postJSON(t, srv, "/api/v1/vehicles/"+created.ID+"/stock-out", `{"stocked_out_on":"2026-10-27"}`), http.StatusBadRequest, "VALIDATION_ERROR")

	viaPatch := patchJSON(t, srv, "/api/v1/vehicles/"+created.ID, `{"stocked_out_on":"2026-10-01"}`)
	assertAPIError(t, viaPatch, http.StatusBadRequest, "VALIDATION_ERROR")

	stayed := decodeVehicle(t, request(t, srv, http.MethodGet, "/api/v1/vehicles/"+created.ID).Body.Bytes()).Data
	if stayed.StockedOutOn == nil || *stayed.StockedOutOn != "2026-10-20" || stayed.StockState != "stocked_out" {
		t.Fatalf("after rejected changes = %+v, want stocked_out 2026-10-20", stayed)
	}
}

func recordVehicle(t *testing.T, srv handler, vin, modelID, enteredOn string) vehicleJSON {
	t.Helper()
	rec := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(vin, modelID, enteredOn))
	if rec.Code != http.StatusCreated {
		t.Fatalf("record status = %d, body = %s", rec.Code, rec.Body.String())
	}
	return decodeVehicle(t, rec.Body.Bytes()).Data
}

func stockOutVehicle(t *testing.T, srv handler, id, stockedOutOn string) vehicleJSON {
	t.Helper()
	rec := postJSON(t, srv, "/api/v1/vehicles/"+id+"/stock-out", `{"stocked_out_on":"`+stockedOutOn+`"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("stock-out status = %d, body = %s", rec.Code, rec.Body.String())
	}
	return decodeVehicle(t, rec.Body.Bytes()).Data
}

func vehicleBodyJSON(vin, modelID, enteredOn string) string {
	raw, err := json.Marshal(map[string]string{
		"vin":        vin,
		"model_id":   modelID,
		"entered_on": enteredOn,
	})
	if err != nil {
		panic(err)
	}
	return string(raw)
}

type errorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ErrorCode string `json:"error_code"`
	} `json:"data"`
}

func assertAPIError(t *testing.T, rec *httptest.ResponseRecorder, status int, errorCode string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("status = %d, body = %s, want %d", rec.Code, rec.Body.String(), status)
	}
	var body errorJSON
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("error body: %v", err)
	}
	if body.Code != status || body.Data.ErrorCode != errorCode {
		t.Fatalf("error = %+v, want code %d %s", body, status, errorCode)
	}
}

type listJSON struct {
	Meta struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	} `json:"meta"`
	Data    []vehicleJSON `json:"data"`
	RawData string
}

func listVehicles(t *testing.T, srv handler, query string) listJSON {
	t.Helper()
	rec := request(t, srv, http.MethodGet, "/api/v1/vehicles"+query)
	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var parsed struct {
		Data json.RawMessage `json:"data"`
		Meta struct {
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatal(err)
	}
	var list listJSON
	list.Meta = parsed.Meta
	list.RawData = string(parsed.Data)
	if err := json.Unmarshal(parsed.Data, &list.Data); err != nil {
		t.Fatal(err)
	}
	return list
}
