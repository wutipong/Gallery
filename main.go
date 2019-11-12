package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"flag"
	"log"
)

func main() {
	address := flag.String("address", ":6969", "Specify the server address")
	path := flag.String("path", "C:\\Users\\mr_ta\\Desktop\\testRepo", "Specifiy the image source path")

	BaseDirectory = *path

	log.Printf("Image Source Path: %s", *path)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/listdir", ListDir)
	e.GET("/listdir/*", ListDir)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
