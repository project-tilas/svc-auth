package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type (
	health struct {
		ServiceName string `json:"serviceName"`
		Alive       bool   `json:"alive"`
		Version     string `json:"version"`
		PodName     string `json:"podName"`
		NodeName    string `json:"nodeName"`
	}
)

var version string
var addr string
var dbAddr string
var nodeName string
var podName string

func init() {
	fmt.Println("Running SVC_AUTH version: " + version)
	addr = getEnvVar("SVC_AUTH_ADDR", ":8080")
	dbAddr = getEnvVar("SVC_AUTH_DB_ADDR", ":8080")
	nodeName = getEnvVar("SVC_AUTH_NODE_NAME", "N/A")
	podName = getEnvVar("SVC_AUTH_POD_NAME", "N/A")
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/health", func(c echo.Context) error {
		u := health{
			Alive:       true,
			ServiceName: "svc-auth",
			Version:     version,
			PodName:     podName,
			NodeName:    nodeName,
		}
		return c.JSON(http.StatusOK, u)
	})

	// Start server
	e.Logger.Fatal(e.Start(addr))
}

func getEnvVar(env string, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
