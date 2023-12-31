package main

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	g := e.Group("/demo/api/v1")

	g.GET("/", func(c echo.Context) error {
		slog.Info("Hello World",
			slog.String("key1", "value1"),
		)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	})

	g.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
		})
	})

	g.GET("/hello/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello " + name,
		})
	})

	g.GET("/subdomain", func(c echo.Context) error {
		host := c.Request().Host
		return c.JSON(http.StatusOK, map[string]string{
			"subdomain": getSubdomain(host),
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func getSubdomain(host string) string {

	host = strings.TrimSpace(host)

	hostParts := strings.Split(host, ".")
	parts := len(hostParts)
	if parts > 2 {
		subdomain := strings.Join(hostParts[:parts-2], ".")
		if subdomain == "www" {
			return ""
		}
		return strings.TrimPrefix(subdomain, "www.")
	}

	return ""
}
