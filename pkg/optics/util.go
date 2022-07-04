package optics

import (
	"io/fs"
	"net/url"
	"path/filepath"
	"strings"
)

// Check working directory for .env files.
func envExists() (exists bool) {
	filepath.WalkDir(".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == "env" {
			exists = true
		} else {
			exists = false
		}
		return nil
	})
	return
}

// Get a short-name for the URL.
func epName(url url.URL) string {
	str := strings.Split(url.String(), "/")
	s := str[len(str)-1]
	return s
}
