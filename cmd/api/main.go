// Package main is the entry point for the application.
package main

import (
	"golang-connect-marketplace/internal/auth/http/handlers"
	"golang-connect-marketplace/internal/auth/http/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	authHandler := handlers.NewAuthHandler()
	routes.RegisterRoutes(e, authHandler)

	e.Start(":6767")
}
