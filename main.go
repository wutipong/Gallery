package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namsral/flag"
)

func main() {
	address := flag.String("address", ":80", "The server address")
	path := flag.String("data_path", "./data", "Image source path")
	prefix := flag.String("url_prefix", "*", "Url prefix")

	flag.Parse()

	BaseDirectory = *path

	log.Printf("Image Source Path: %s", *path)
	log.Printf("using prefix %s", *prefix)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Rewrite(map[string]string{
		*prefix: "$1",
	}))

	// Routes
	e.GET("/", hello)
	e.GET("/browse", browse)
	e.GET("/browse/*", browse)

	e.GET("/view", view)
	e.GET("/view/*", view)

	e.Static("/static", "static")

	e.GET("/get_image/*", GetImage)

	e.GET("/get_cover", GetCover)
	e.GET("/get_cover/*", GetCover)

	e.GET("/view", view)
	e.GET("/view/*", view)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}

// Handler
func hello(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/browse")
}
