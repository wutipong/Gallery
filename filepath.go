package main

import (
	"os"
	"strings"
)

//BaseDirectory the data directory
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
func ListDir(p string) (dirs []FileEntry, files []FileEntry, err error) {

	fullpath := BaseDirectory + string(os.PathSeparator) + p

	dir, err := os.Open(fullpath)
	if err != nil {
		return
	}
	children, err := dir.Readdir(0)
	if err != nil {
		return
	}

	for _, f := range children {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}

		if !f.IsDir() && !filter(f.Name()) {
			continue
		}

		if f.IsDir() {
			dirs = append(dirs, FileEntry{
				Filename: f.Name(),
				IsDir:    f.IsDir(),
			})
		} else {
			files = append(files, FileEntry{
				Filename: f.Name(),
				IsDir:    f.IsDir(),
			})
		}
	}

	return
}
