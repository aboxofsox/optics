package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
)

func Reverse(origin string) {
	u, err := url.Parse(origin)
	if err != nil {
		log.Fatal(err)
	}

	proxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = u.Scheme
			r.URL.Host = u.Host
			r.Host = u.Host

			r.URL.Path = u.Path + strings.TrimRight(r.URL.Path, "/") + "/"

			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		},
	}

	fep := httptest.NewServer(proxy)
	defer fep.Close()

	res, err := http.Get(fep.URL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}
