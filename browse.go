package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func init() {
	var err error
	folderTemplate, err = template.New("folderitem.gohtml").ParseFiles("template/folderitem.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}

	fileTemplate, err = template.New("fileitem.gohtml").ParseFiles("template/fileitem.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}

	navTemplate, err = template.New("browse-nav.gohtml").ParseFiles("template/browse-nav.gohtml")
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
}

var folderTemplate *template.Template = nil
var fileTemplate *template.Template = nil
var navTemplate *template.Template = nil

type folderItem struct {
	Name     string
	LinkURL  string
	ThumbURL string
}

type fileItem struct {
	Name     string
	LinkURL  string
	ThumbURL string
}

type navItem struct {
	Name string
	URL  string
}

//WriteBreadcrumb write breadcrumb component.
func WriteBreadcrumb(writer io.Writer, path string) {
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

	navTemplate.Execute(writer, struct{ NavItems []navItem }{NavItems: items})
}

// WriteDirectories write directory entries.
func WriteDirectories(writer io.Writer, path string, dirs []FileEntry) {
	io.WriteString(writer, `<div class="container">`)
	length := len(dirs)
	for i := 0; i < length; i++ {
		dir := dirs[i]
		if i%3 == 0 {
			io.WriteString(writer, `<div class="row">`)
		}
		var url string
		var thumbURL string
		if path == "" {
			url = "/browse/" + dir.Filename
			thumbURL = "/get_cover/" + dir.Filename
		} else {
			url = "/browse/" + path + "/" + dir.Filename
			thumbURL = "/get_cover/" + path + "/" + dir.Filename
		}

		folderTemplate.Execute(writer, folderItem{Name: dir.Filename, LinkURL: url, ThumbURL: thumbURL})

		if i%3 == 2 || i == length-1 {
			io.WriteString(writer, `</div>`)
		}
	}
	io.WriteString(writer, `</div>`)
}

// WriteFiles write file entries.
func WriteFiles(writer io.Writer, path string, files []FileEntry) {
	io.WriteString(writer, `<div class="container">`)
	length := len(files)
	for i := 0; i < length; i++ {
		file := files[i]
		if i%3 == 0 {
			io.WriteString(writer, `<div class="row">`)
		}
		var url string
		if path == "" {
			url = "/get_image/" + file.Filename
		} else {
			url = "/get_image/" + path + "/" + file.Filename
		}

		//io.WriteString(writer, fmt.Sprintf(`<div class="col"><a href="%s">%s</a></div>`, url, file.Filename))
		fileTemplate.Execute(writer, fileItem{Name: file.Filename, LinkURL: url, ThumbURL: url})

		if i%3 == 2 || i == length-1 {
			io.WriteString(writer, `</div>`)
		}
	}
	io.WriteString(writer, `</div>`)
}

// Handler
func browse(c echo.Context) error {
	builder := strings.Builder{}

	WriteHeader(&builder, Header{Title: "Hello"})

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	if strings.HasSuffix(p, "/") {
		p = p[0 : len(p)-1]
		return c.Redirect(http.StatusPermanentRedirect, "/browse/"+p)
	}

	dirs, files, err := ListDir(p)
	if err != nil {
		return err
	}

	WriteBreadcrumb(&builder, p)
	WriteDirectories(&builder, p, dirs)
	WriteFiles(&builder, p, files)

	return c.HTML(http.StatusOK, builder.String())
}
