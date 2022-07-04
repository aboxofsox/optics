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
	"github.com/aboxofsox/optics/pkg/logger"
)

/*
Proxy GET requests.

This function creates a temporary HTTP server using httptest.NewServer().
*/
func (ctrl *Controller) Proxy(origin string, done func()) {
	defer done()
	var resMsg string
	var resStatusCode string

	lgr := logger.New()
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
		resMsg = colors.Red(http.StatusText(res.StatusCode))
	} else {
		resStatusCode = colors.Green(res.StatusCode)
		resMsg = colors.Green(http.StatusText(res.StatusCode))
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("unable to read response body: %s\n", err.Error())
		return
	}

	li := &logger.LogItem{
		Timestamp:         time.Now().Format(time.ANSIC),
		Endpoint:          "[PROXIED]" + origin,
		StatusCode:        res.StatusCode,
		StatusCodeMessage: strings.ToUpper(http.StatusText(res.StatusCode)),
		Elapsed:           since,
	}
	lgr.Stash(li)
	lgr.Write("./res/api.log")

	ctrl.Json(b)
	fmt.Printf(
		"[PROXIED] %s: %s %s - %.2dms\n",
		u.String(),
		resStatusCode,
		resMsg,
		since.Milliseconds(),
	)

}
