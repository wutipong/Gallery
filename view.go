package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"

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

type viewData struct {
	Name       string
	Title      string
	BrowseURL  string
	ImageURLs  []string
	StartIndex int64
}

func view(c echo.Context) error {
	builder := strings.Builder{}

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	_, files, err := ListDir(p)
	if err != nil {
		return err
	}

	var startIndex int64 = 1
	if i, e := strconv.ParseInt(c.QueryParam("index"), 10, 64); e == nil {
		startIndex = i
	}

	data := viewData{
		Name: 		p,
		Title:      fmt.Sprintf("Gallery - Viewing [%s]", p),
		BrowseURL:  "/browse/" + p,
		ImageURLs:  createImageURLs(p, files),
		StartIndex: startIndex,
	}
	err = viewTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}

func createImageURLs(path string, files []FileEntry) []string {
	output := make([]string, len(files))
	for i, file := range files {
		var url string
		if path == "" {
			url = "/get_image/" + file.Filename
		} else {
			url = "/get_image/" + path + "/" + file.Filename
		}

		output[i] = url
	}
	return output
}
