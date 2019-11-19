package main

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

func findCover(files []os.FileInfo) os.FileInfo {
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := strings.ToLower(f.Name())
		if strings.Contains(name, "cover") {
			return f
		}
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !filter(f.Name()) {
			continue
		}
		return f
	}

	return nil
}

// GetCover returns a cover image with specific width/height while retains aspect ratio.
func GetCover(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}
	fullpath := BaseDirectory + string(os.PathSeparator) + p

	var width, height int64 = 0, 0
	if w, e := strconv.ParseInt(c.QueryParam("width"), 10, 64); e == nil {
		width = w
	}

	if h, e := strconv.ParseInt(c.QueryParam("width"), 10, 64); e == nil {
		height = h
	}

	dir, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	children, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	cover := findCover(children)

	if cover == nil {
		return c.File("static/img/notfound_thumb.png")
	}

	coverPath := fullpath + string(os.PathSeparator) + cover.Name()

	if width == 0 || height == 0 {
		return c.File(coverPath)
	}

	file, err := os.Open(coverPath)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(file)

	if err != nil {
		return err
	}

	resized := resize.Thumbnail(uint(width), uint(height), img, resize.MitchellNetravali)
	buffer := bytes.Buffer{}

	jpeg.Encode(&buffer, resized, nil)
	return c.Blob(http.StatusOK, "image/jpeg", buffer.Bytes())
}
