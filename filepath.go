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

// FileEntry is an entry of returning record.
type FileEntry struct {
	Filename string `json:"filename"`
	IsDir    bool   `json:"is_dir"`
}

// ListDir returns a list of content of a directory.
func ListDir(c echo.Context) error {
	p, err := url.PathUnescape(c.Param("*"))
	if err != nil {
		return err
	}

	fullpath := BaseDirectory + string(os.PathSeparator) + p

	dir, err := os.Open(fullpath)
	if err != nil {
		return err
	}
	children, err := dir.Readdir(0)
	if err != nil {
		return err
	}

	output := []FileEntry{}

	for _, f := range children {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if !f.IsDir() && !filter(f.Name()) {
			continue
		}

		output = append(output, FileEntry{
			Filename: f.Name(),
			IsDir:    f.IsDir(),
		})
	}

	return c.JSON(http.StatusOK, output)
}
