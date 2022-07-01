package optics

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/aboxofsox/optics/pkg/colors"
)

/*
Proxy GET requests.

This function creates a temporary HTTP server using httptest.NewServer().
*/
func (ctrl *Controller) Proxy(origin string, done func()) {
	defer done()
	var resMsg string
	var resStatusCode string
	u, err := url.Parse(origin)
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

	start := time.Now()
	res, err := http.Get(fep.URL)
	if err != nil {
		fmt.Printf("unable to do get request: %s\n", err.Error())
		return
	}
	since := time.Since(start)

	if res.StatusCode == http.StatusNotFound {
		resStatusCode = colors.Red(res.StatusCode)
		resMsg = colors.Red(StatusCodes[res.StatusCode])
	} else {
		resStatusCode = colors.Green(res.StatusCode)
		resMsg = colors.Green(StatusCodes[res.StatusCode])
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to read response body: %s\n", err.Error())
		return
	}

	ctrl.Json(b)
	ctrl.Log(*res, since.Seconds())
	fmt.Printf(
		"[PROXIED] %s: %s %s - %v\n",
		u.String(),
		resStatusCode,
		resMsg,
		colors.Cyan(since),
	)

}
