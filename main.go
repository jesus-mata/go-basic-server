package main

import (
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		slog.Info("Hello World",
			slog.String("key1", "value1"),
		)
		return c.JSON(200, map[string]string{
			"message": "Hello World",
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "OK",
		})
	})

	e.GET("/subdomain", func(c echo.Context) error {
		host := c.Request().Host
		return c.JSON(200, map[string]string{
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
