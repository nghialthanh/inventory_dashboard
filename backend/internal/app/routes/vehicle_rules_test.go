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
	"inventory_dashboard/internal/repositories"
	"inventory_dashboard/internal/services"
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

	srv := routes.New(routes.Config{
		Services: services.New(services.Deps{
			Catalog:  repositories.NewCatalog(pool),
			Vehicles: repositories.NewVehicles(pool),
			Today:    func() time.Time { return today },
		}),
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

func recordVehicle(t *testing.T, srv handler, vin, modelID, enteredOn string) vehicleJSON {
	t.Helper()
	rec := postJSON(t, srv, "/api/v1/vehicles", vehicleBodyJSON(vin, modelID, enteredOn))
	if rec.Code != http.StatusCreated {
		t.Fatalf("record status = %d, body = %s", rec.Code, rec.Body.String())
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
