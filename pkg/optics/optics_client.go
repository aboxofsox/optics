package optics

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Do(url url.URL) http.Response {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer res.Body.Close()

	return *res
}

func Reader(r io.Reader) []byte {
	bd, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err.Error())
	}
	return bd
}

func Json(b []byte) interface{} {
	var obj any
	if err := json.Unmarshal(b, &obj); err != nil {
		log.Fatal(err.Error())
	}
	return obj
}




