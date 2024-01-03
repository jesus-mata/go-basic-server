package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {

	if len(os.Args) <= 1 {
		runServer()
	}

	action := os.Args[1]

	switch action {
	case "healthcheck":
		healthCheck("/demo/api/v1/health")
	case "run":
		runServer()
	default:
		fmt.Println("Usage: ./main <action>")
		fmt.Println("action: run | healthcheck")
		fmt.Println("Example: ./main run or ./main healthcheck")
		fmt.Println("If no action is provided, then the application will run")
		panic("Unknown action")
	}
}

func runServer() {
	appName := os.Getenv("APP_NAME")

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

		hostname := getHostname()

		slog.Info("Request to hello endpoint",
			slog.String("name", name),
			slog.String("host", c.Request().Host),
			slog.String("hostname", getHostname()),
		)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello " + name + " from " + hostname + "!",
		})
	})

	g.GET("/info", func(c echo.Context) error {
		hostname := getHostname()
		return c.JSON(http.StatusOK, map[string]string{
			"host": hostname,
			"app":  appName,
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

// Returns the hostname of the machine
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error getting hostname",
			slog.String("error", err.Error()),
		)
		hostname = "unknown"
	}
	return hostname
}

// Check the health of the application
// Make a request to the health endpoint and check the status code
// If the status code is 200, then the application is healthy and return exit code 0
// If the status code is not 200, then the application is not healthy and return exit code 1
func healthCheck(healthEndpoint string) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	url := "http://localhost:8080/" + strings.TrimPrefix(healthEndpoint, "/")
	//demo/api/v1/health
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Status Code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}

}
