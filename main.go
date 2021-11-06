package main

import (
	"log"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namsral/flag"
	urlutil "github.com/wutipong/go-utils/url"
)

func main() {
	address := flag.String("address", ":80", "The server address")
	dataPath := flag.String("data_path", "./data", "Image source path")
	prefix := flag.String("url_prefix", "", "Url prefix")

	flag.Parse()

	BaseDirectory = *dataPath

	log.Printf("Image Source Path: %s", *dataPath)
	log.Printf("using prefix %s", *prefix)
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Pre(middleware.RemoveTrailingSlash())
	if *prefix != "" {
		pattern := path.Join(*prefix, "*")
		e.Pre(middleware.Rewrite(map[string]string{
			*prefix: "/",
			pattern: "/$1",
		}))

		urlutil.SetPrefix(*prefix)
	}

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
