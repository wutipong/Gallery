package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").
		ParseFiles(
			"template/view.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template = nil

func view(c echo.Context) error {
	builder := strings.Builder{}

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	dirs, files, err := ListDir(p)
	if err != nil {
		return err
	}
	data := browseData{
		Title:    fmt.Sprintf("Gallery - Viewing [%s]", p),
		NavItems: createBreadcrumb(p),
		Files:    createFileItems(p, files),
		Dirs:     createDirectoryItems(p, dirs),
	}
	err = viewTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
