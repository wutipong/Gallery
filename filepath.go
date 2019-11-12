package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

var BaseDirectory string
var filter func(path string) bool

func init() {
	filter = func(path string) bool {
		nameLower := strings.ToLower(path)
		if strings.HasSuffix(nameLower, ".jpeg") {
			return true
		}
		if strings.HasSuffix(nameLower, ".jpg") {
			return true
		}
		if strings.HasSuffix(nameLower, ".png") {
			return true
		}
		return false
	}
}

type FileEntry struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"is_dir"`
}

func ListDir(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	fullpath := BaseDirectory + string(os.PathSeparator) + p

	file, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	dirs, err := file.Readdir(0)
	if err != nil {
		return err
	}

	output := []FileEntry{}

	for _, dir := range dirs {
		if strings.HasPrefix(dir.Name(), ".") {
			continue
		}

		output = append(output, FileEntry{
			Filename: dir.Name(),
			IsDir:    dir.IsDir(),
		})
	}

	return c.JSON(http.StatusOK, output)
}
