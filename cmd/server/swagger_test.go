package main

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func TestSwaggerUIAndSpec(t *testing.T) {
	app := fiber.New()
	app.Use(swagger.New(swagger.Config{FilePath: "../../docs/openapi.yaml", Path: "swagger"}))
	for _, path := range []string{"/swagger", "/docs/openapi.yaml"} {
		resp, err := app.Test(httptest.NewRequest("GET", path, nil))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Errorf("GET %s: status %d", path, resp.StatusCode)
		}
	}
}
