package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"inventory_dashboard/internal/app/routes"
)

func TestHealthRespondsOutsideTheVersionedAPI(t *testing.T) {
	srv := routes.New(routes.Config{})

	rec := request(t, srv, http.MethodGet, "/health")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /health status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("GET /health content-type = %q, want application/json", ct)
	}

	var body struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("GET /health body: %v", err)
	}
	if body.Code != 200 || body.Message != "ok" {
		t.Fatalf("GET /health envelope = code %d message %q, want 200 ok", body.Code, body.Message)
	}
	if string(body.Data) != "{}" {
		t.Fatalf("GET /health data = %s, want {}", body.Data)
	}

	versioned := request(t, srv, http.MethodGet, "/api/v1/health")
	if versioned.Code == http.StatusOK {
		t.Fatal("GET /api/v1/health returned 200; health belongs on /health, outside /api/v1")
	}
}

func TestMetricsExposesPrometheusText(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>Inventory</title>"), 0o644); err != nil {
		t.Fatal(err)
	}
	srv := routes.New(routes.Config{StaticDir: dir})

	rec := request(t, srv, http.MethodGet, "/metrics")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want 200", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/plain") && !strings.Contains(ct, "openmetrics") {
		t.Fatalf("GET /metrics content-type = %q, want Prometheus text", ct)
	}
	body := rec.Body.String()
	if strings.Contains(body, "<title>Inventory</title>") {
		t.Fatal("GET /metrics fell back to the inventory shell")
	}
	if !strings.Contains(body, "http_requests_total") && !strings.Contains(body, "# HELP") {
		t.Fatalf("GET /metrics body missing Prometheus exposition, got %q", body[:min(120, len(body))])
	}

	versioned := request(t, srv, http.MethodGet, "/api/v1/metrics")
	if versioned.Code == http.StatusOK && strings.Contains(versioned.Body.String(), "http_requests_total") {
		t.Fatal("GET /api/v1/metrics returned Prometheus text; metrics belongs on /metrics, outside /api/v1")
	}
}

func TestDashboardAndAPIShareOneOrigin(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<!doctype html><title>Inventory</title>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}

	srv := routes.New(routes.Config{StaticDir: dir})

	page := request(t, srv, http.MethodGet, "/")
	if page.Code != http.StatusOK || !strings.Contains(page.Body.String(), "<title>Inventory</title>") {
		t.Fatalf("GET / status = %d body = %q, want the inventory shell", page.Code, page.Body.String())
	}

	asset := request(t, srv, http.MethodGet, "/app.js")
	if asset.Code != http.StatusOK || asset.Body.String() != "console.log(1)" {
		t.Fatalf("GET /app.js status = %d body = %q, want the built file", asset.Code, asset.Body.String())
	}

	fallback := request(t, srv, http.MethodGet, "/inventory")
	if fallback.Code != http.StatusOK || !strings.Contains(fallback.Body.String(), "<title>Inventory</title>") {
		t.Fatalf("GET /inventory status = %d body = %q, want the shell fallback", fallback.Code, fallback.Body.String())
	}

	health := request(t, srv, http.MethodGet, "/health")
	if health.Code != http.StatusOK || !strings.Contains(health.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("GET /health status = %d content-type = %q, want JSON", health.Code, health.Header().Get("Content-Type"))
	}

	api := request(t, srv, http.MethodGet, "/api/v1/vehicles")
	if strings.Contains(api.Body.String(), "<title>Inventory</title>") {
		t.Fatal("an /api path fell back to the inventory shell")
	}
}

type handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func request(t *testing.T, srv handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}
