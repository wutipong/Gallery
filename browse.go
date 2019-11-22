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
	broseTemplate, err = template.New("browse.gohtml").
		ParseFiles(
			"template/browse.gohtml",
			"template/header.gohtml",
		)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var broseTemplate *template.Template = nil

type browseData struct {
	Title    string
	NavItems []navItem
	Files    []fileItem
	Dirs     []folderItem
}

type folderItem struct {
	Name     string
	LinkURL  string
	ThumbURL string
}

type fileItem struct {
	Name       string
	Path       string
	ImageURL   string
	StartIndex int
}

type navItem struct {
	Name string
	URL  string
}

func createBreadcrumb(path string) []navItem {
	items := []navItem{}

	items = append(items, navItem{
		Name: "Home",
		URL:  "/browse",
	})
	if path != "" {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			items = append(items, navItem{
				Name: part,
				URL:  "/browse/" + PathLevel(path, i+1),
			})
		}
	}
	return items
}

func createDirectoryItems(path string, dirs []FileEntry) []folderItem {
	output := make([]folderItem, len(dirs))
	for i, dir := range dirs {
		var url string
		var thumbURL string
		if path == "" {
			url = "/browse/" + dir.Filename
			thumbURL = "/get_cover/" + dir.Filename
		} else {
			url = "/browse/" + path + "/" + dir.Filename
			thumbURL = "/get_cover/" + path + "/" + dir.Filename
		}

		output[i] = folderItem{Name: dir.Filename, LinkURL: url, ThumbURL: thumbURL}
	}
	return output
}

func createFileItems(path string, files []FileEntry) []fileItem {
	output := make([]fileItem, len(files))
	for i, file := range files {
		var url string
		if path == "" {
			url = file.Filename
		} else {
			url = path + "/" + file.Filename
		}

		output[i] = fileItem{Name: file.Filename, Path: path, ImageURL: url, StartIndex: i + 1}
	}
	return output
}

// Handler
func browse(c echo.Context) error {
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
		Title:    fmt.Sprintf("Gallery - Browsing [%s]", p),
		NavItems: createBreadcrumb(p),
		Files:    createFileItems(p, files),
		Dirs:     createDirectoryItems(p, dirs),
	}
	err = broseTemplate.Execute(&builder, data)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.HTML(http.StatusOK, builder.String())
}
