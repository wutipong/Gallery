package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

//WriteBreadcrumb write breadcrumb component.
func WriteBreadcrumb(writer io.Writer, path string) {
	io.WriteString(writer, "<nav>")
	io.WriteString(writer, `<ol class="breadcrumb">`)
	if path == "" {
		io.WriteString(writer, `<li class="breadcrumb-item active" aria-current="page">Home</li>`)
	} else {
		io.WriteString(writer, `<li class="breadcrumb-item active" aria-current="page"><a href="/browse">Home</a></li>`)

		parts := strings.Split(path, "/")
		for i, part := range parts {
			if i == len(parts)-1 {
				io.WriteString(writer, fmt.Sprintf(`<li class="breadcrumb-item active" aria-current="page">%s</li>`, part))
			} else {
				url := PathLevel(path, i+1)
				url = "/browse/" + url
				io.WriteString(writer, fmt.Sprintf(`<li class="breadcrumb-item active" aria-current="page"><a href="%s">%s</a></li>`, url, part))
			}
		}
	}
	io.WriteString(writer, "</ol>")
	io.WriteString(writer, "</nav>")
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

		if path == "" {
			url = "/browse/" + dir.Filename
		} else {
			url = "/browse/" + path + "/" + dir.Filename
		}

		io.WriteString(writer, fmt.Sprintf(`<div class="col"><a href="%s">%s</a></div>`, url, dir.Filename))

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

		io.WriteString(writer, fmt.Sprintf(`<div class="col"><a href="%s">%s</a></div>`, url, file.Filename))

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
