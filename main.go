package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"flag"
	"log"
)

func main() {
	address := flag.String("address", ":6969", "Specify the server address")
	path := flag.String("path", "/data", "Specifiy the image source path")

	BaseDirectory = *path

	log.Printf("Image Source Path: %s", *path)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/browse", browse)
	e.GET("/browse/*", browse)

	e.Static("/static", "static")

	e.GET("/get_image/*", GetImage)

	e.GET("/get_cover", GetCover)
	e.GET("/get_cover/*", GetCover)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handler
func hello(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/browse")
}
