package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"inventory_dashboard/api/middlewares"
	"inventory_dashboard/internal/common/trace"
)

func TestTraceGeneratesAndPropagates(t *testing.T) {
	e := echo.New()
	e.Use(middlewares.Trace())
	var got string
	e.GET("/ping", func(c echo.Context) error {
		got = trace.ID(c.Request().Context())
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if got == "" {
		t.Fatal("expected trace_id on context")
	}
	if rec.Header().Get(trace.Header) != got {
		t.Fatalf("response header %q, want %q", rec.Header().Get(trace.Header), got)
	}
}

func TestTraceAcceptsIncomingHeader(t *testing.T) {
	e := echo.New()
	e.Use(middlewares.Trace())
	var got string
	e.GET("/ping", func(c echo.Context) error {
		got = trace.ID(c.Request().Context())
		return c.NoContent(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(trace.Header, "client-trace-1")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if got != "client-trace-1" {
		t.Fatalf("trace_id = %q, want client-trace-1", got)
	}
}
