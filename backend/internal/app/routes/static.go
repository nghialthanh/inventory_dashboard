package routes

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func mountStatic(e *echo.Echo, dir string) {
	if dir == "" {
		return
	}
	serve := func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		if urlPath == "/api" || strings.HasPrefix(urlPath, "/api/") {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		return c.File(resolveStatic(dir, urlPath))
	}
	e.GET("/", serve)
	e.GET("/*", serve)
}

func resolveStatic(dir, urlPath string) string {
	index := filepath.Join(dir, "index.html")
	cleaned := path.Clean("/" + urlPath)
	rel := strings.TrimPrefix(cleaned, "/")
	if rel == "" || rel == "." {
		return index
	}
	full := filepath.Join(dir, filepath.FromSlash(rel))
	relToDir, err := filepath.Rel(dir, full)
	if err != nil || relToDir == ".." || strings.HasPrefix(relToDir, ".."+string(os.PathSeparator)) {
		return index
	}
	info, err := os.Stat(full)
	if err != nil || info.IsDir() {
		return index
	}
	return full
}
