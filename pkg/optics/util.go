package optics

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
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
func cPath(args ...string) string {
	w := strings.Builder{}
	for i, arg := range args {
		if arg == "/" || arg == "\\" {
			args[i] = string(filepath.Separator)
		}
		w.Write([]byte(arg))
	}
	return w.String()
}

func read(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read file: %v", err)
	}

	return data, nil
}

func parse(data []byte, envMap map[string]string) (string, error) {
	w := &strings.Builder{}
	tmpl, err := template.New("toml").Parse(string(data))
	if err != nil {
		return "", fmt.Errorf("parse config: %v", err)
	}

	tmpl.Execute(w, envMap)

	return w.String(), nil
}

func decode(conf string) (*Config, error) {
	var config *Config
	if _, err := toml.Decode(conf, &config); err != nil {
		return nil, fmt.Errorf("decode template: %v", err)
	}
	return config, nil
}

func opTime(start time.Time) time.Duration { return time.Since(start) }

func logWriter(p string, data []byte) (int, error) {
	file, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, PERMS)
	if err != nil {
		return 0, fmt.Errorf("log writer: %v", err)
	}
	defer file.Close()
	n, err := file.Write(data)
	if err != nil {
		return 0, fmt.Errorf("write file: %v", err)
	}
	return n, nil
}

func timestamp(t time.Time) string { return t.Format(time.ANSIC) }
