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
				url := strings.Join(parts[0:i+1], "/")
				url = "/browse/" + url
				io.WriteString(writer, fmt.Sprintf(`<li class="breadcrumb-item active" aria-current="page"><a href="%s">%s</a></li>`, url, part))
			}
		}
	}
	io.WriteString(writer, "</ol>")
	io.WriteString(writer, "</nav>")
}

// Handler
func browse(c echo.Context) error {
	builder := strings.Builder{}

	WriteHeader(&builder, Header{Title: "Hello"})

	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	entries, err := ListDir(p)
	if err != nil {
		return err
	}

	entries[0].Filename = "a"
	WriteBreadcrumb(&builder, p)
	builder.WriteString("<p>Hello World</p>")

	return c.HTML(http.StatusOK, builder.String())
}
