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
	urlutil "github.com/wutipong/go-utils/url"
)

func init() {
	var err error
	broseTemplate, err = template.New("browse.gohtml").
		Funcs(urlutil.HtmlTemplateFuncMap()).
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
	Dirs     []fileItem
}

type fileItem struct {
	Name     string
	ThumbURL string
	LinkURL  string
}

type navItem struct {
	Name string
	URL  string
}

func createBreadcrumb(path string) []navItem {
	items := []navItem{}

	items = append(items, navItem{
		Name: "Home",
		URL:  urlutil.CreateURL("/browse"),
	})
	if path != "" {
		parts := strings.Split(path, "/")
		for i, part := range parts {
			items = append(items, navItem{
				Name: part,
				URL:  urlutil.CreateURL("/browse/", PathLevel(path, i+1)),
			})
		}
	}
	return items
}

func createDirectoryItems(path string, dirs []FileEntry) []fileItem {
	output := make([]fileItem, len(dirs))
	for i, dir := range dirs {
		var url string
		var thumbURL string
		if path == "" {
			url = urlutil.CreateURL("/browse/", dir.Filename)
			thumbURL = urlutil.CreateURL("/get_cover/", dir.Filename)
		} else {
			url = urlutil.CreateURL("/browse/", path, dir.Filename)
			thumbURL = urlutil.CreateURL("/get_cover/", path, dir.Filename)
		}

		output[i] = fileItem{Name: dir.Filename, LinkURL: url, ThumbURL: thumbURL}
	}
	return output
}

func createFileItems(path string, files []FileEntry) []fileItem {
	output := make([]fileItem, len(files))
	viewURL := urlutil.CreateURL("view", path)
	for i, file := range files {
		var url string
		if path == "" {
			url = urlutil.CreateURL("get_image", file.Filename)
		} else {
			url = urlutil.CreateURL("get_image", path, file.Filename)
		}

		index := i + 1
		output[i] = fileItem{
			Name:     file.Filename,
			ThumbURL: url,
			LinkURL:  fmt.Sprintf("%s#%d", viewURL, index),
		}
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
