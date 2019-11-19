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

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

// GetImage returns an image with specific width/height while retains aspect ratio.
func GetImage(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}
	fullpath := BaseDirectory + string(os.PathSeparator) + p

	var width, height int64 = 0, 0
	if w, e := strconv.ParseInt(c.QueryParam("width"), 10, 64); e == nil {
		width = w
		height = width
	}

	if h, e := strconv.ParseInt(c.QueryParam("height"), 10, 64); e == nil {
		height = h
	}

	if width == 0 || height == 0 {
		return c.File(fullpath)
	}

	file, err := os.Open(fullpath)
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
