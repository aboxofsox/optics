package optics

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Create a new HTTP client controller
func New() *Controller {
	env := readEnv()
	p := cPath(".", string(filepath.Separator), "optics.toml")

	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err.Error())
	}
	d, err := read(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	t, err := parse(d, env)
	if err != nil {
		log.Fatal(err.Error())
	}
	dc, err := decode(t)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &Controller{
		Config: dc,
		Buffer: &bytes.Buffer{},
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		HttpResponse: &HttpResponse{},
	}
}
