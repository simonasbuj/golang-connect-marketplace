// Package main is the entry point for the application.
package main

import (
	"golang-connect-marketplace/internal/auth/http/handlers"
	"golang-connect-marketplace/internal/auth/http/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	authHandler := handlers.NewHandler()
	routes.RegisterRoutes(e, authHandler)

	err := e.Start(":6767")
	if err != nil {
		log.Error(err)
	}
}
