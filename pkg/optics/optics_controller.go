package optics

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// Create a new HTTP client controller
func New() *Controller {
	var envMap map[string]string

	if envExists() {
		envMap, _ = godotenv.Read()
	}
	p := cPath(".", string(filepath.Separator), "optics.toml")

	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err.Error())
	}
	d, err := read(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	t, err := parse(d, envMap)
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
