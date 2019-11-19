package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	viewTemplate, err = template.New("view.gohtml").ParseFiles("template/view.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var viewTemplate *template.Template = nil

type viewData struct{}

// Handler
func view(c echo.Context) error {
	builder := strings.Builder{}

	WriteHeader(&builder, Header{Title: "View"})
	/*
		p, err := url.PathUnescape(c.Param("*"))
		if err != nil {
			return err
		}
		_, files, err := ListDir(p)
		if err != nil {
			return err
		}
	*/
	viewTemplate.Execute(&builder, viewData{})

	return c.HTML(http.StatusOK, builder.String())
}
