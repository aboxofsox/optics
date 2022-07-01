package optics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aboxofsox/optics/pkg/colors"
)

type Config struct {
	Name        string            `toml:"name"`
	Scheme      string            `toml:"scheme"`
	Host        string            `toml:"host"`
	Endpoints   []string          `toml:"endpoints"`
	QueryParams map[string]string `toml:"query_params"`
	Headers     map[string]string `tom:"headers"`
	OutFile     bool              `toml:"outfile"`
	Outdir      string            `toml:"outdir"`
}

type HttpResponse struct {
	StatusCode  int
	Method      string
	Message     string
	ContentType string
	Headers     map[string][]string
	Error       error
	Body        []byte
	Time        float64
}

type Options struct {
	UseProxy bool
}

type Controller struct {
	Config       *Config
	Buffer       *bytes.Buffer
	Client       *http.Client
	Url          *url.URL
	HttpResponse *HttpResponse
	Data         []byte
}

const PERMS = 0770

// Initialize the HTTP client controller.
func (ctrl *Controller) Init(options *Options) {
	var wg sync.WaitGroup
	ctrl.Url = &url.URL{
		Scheme: ctrl.Config.Scheme,
		Host:   ctrl.Config.Host,
	}

	if _, err := os.Stat(ctrl.Config.Outdir); os.IsNotExist(err) {
		os.Mkdir(ctrl.Config.Outdir, PERMS)
	}

	if len(ctrl.Config.Endpoints) == 0 {
		log.Fatal("no endpoints to test")
	}

	for i := range ctrl.Config.Endpoints {
		wg.Add(1)

		ctrl.Url.Path = ctrl.Config.Endpoints[i]

		kv := ctrl.Url.Query()

		if len(ctrl.Config.QueryParams) != 0 {
			for k, v := range ctrl.Config.QueryParams {
				kv.Add(k, v)
			}

			ctrl.Url.RawQuery = kv.Encode()
		}

		if options.UseProxy {
			go ctrl.Proxy(ctrl.Url.String(), wg.Done)
		} else {
			go ctrl.Get(ctrl.Url.String(), wg.Done)
		}

		wg.Wait()
	}
}

/*
Do an HTTP GET request and process the data.

The response is saved as a JSON file and
the results of the response are logged
to a log file in the same directory.
The log file will always be appended to
if it already exists.
*/
func (ctrl *Controller) Get(url string, done func()) {
	defer done()
	var resMsg string
	var resStatusCode string

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	start := time.Now()
	res, err := ctrl.Client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		resStatusCode = colors.Red(res.StatusCode)
		resMsg = colors.Red(StatusCodes[res.StatusCode])
	} else {
		resStatusCode = colors.Green(res.StatusCode)
		resMsg = colors.Green(StatusCodes[res.StatusCode])
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	since := opTime(start)

	ctrl.Json(data)
	ctrl.Log(*res, since)
	fmt.Printf(
		"%s: %s %s - %v seconds\n",
		url,
		resStatusCode,
		resMsg,
		colors.Cyan(since))
	ctrl.Buffer.Reset()

}

// Write the response to a JSON file.
func (ctrl *Controller) Json(data []byte) (int, error) {
	buf := new(bytes.Buffer)
	if len(data) == 0 {
		return 0, fmt.Errorf("no data to write: %d", len(data))
	}
	if err := json.Indent(buf, data, "", " "); err != nil {
		return 0, err
	}

	tmp, err := ioutil.TempFile(
		ctrl.Config.Outdir,
		fmt.Sprintf("%s-*.json",
			epName(*ctrl.Url)),
	)
	if err != nil {
		return 0, err
	}
	defer tmp.Close()

	if _, err := tmp.Write(buf.Bytes()); err != nil {
		return 0, err
	}

	return len(data), nil
}

func mkJson(data []byte) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	if len(data) == 0 {
		return nil, fmt.Errorf("no data to write: %d", len(data))
	}
	if err := json.Indent(buf, data, "", " "); err != nil {
		return nil, err
	}
	return buf, nil
}

func mkTemp(p, pattern string) (*os.File, error) {
	tmp, err := ioutil.TempFile(p, pattern)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

func writeTemp(tmp *os.File, buf bytes.Buffer) (int, error) {
	n, err := tmp.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}

func writeResponse(res http.Response, duration float64) *HttpResponse {
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	return &HttpResponse{
		StatusCode:  res.StatusCode,
		Method:      res.Request.Method,
		Message:     StatusCodes[res.StatusCode],
		ContentType: res.Header.Get("content-type"),
		Headers:     res.Header,
		Error:       errors.New("response invalid"),
		Body:        body,
		Time:        duration,
	}
}

// Handle the log
func (ctrl *Controller) Log(res http.Response, duration float64) {
	headers := map[string][]string{}
	for k, v := range res.Request.Header {
		headers[k] = append(headers[k], v...)
	}

	ctrl.HttpResponse.StatusCode = res.StatusCode
	ctrl.HttpResponse.Method = res.Request.Method
	ctrl.HttpResponse.Time = duration
	ctrl.HttpResponse.Headers = headers

	lfp := filepath.Join(
		".", string(filepath.Separator),
		ctrl.Config.Outdir,
		"optics.log",
	)

	tm := time.Now()
	timestamp := fmt.Sprintf(
		"%d-%02d-%02d %02d:%02d",
		tm.Year(),
		tm.Month(),
		tm.Day(),
		tm.Hour(),
		tm.Minute(),
	)

	file, err := os.OpenFile(lfp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, PERMS)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf(
		"%s %s %d %s - %.2f seconds\n",
		ctrl.Url.String(),
		timestamp,
		res.StatusCode,
		strings.ToUpper(StatusCodes[res.StatusCode]),
		duration,
	)); err != nil {
		log.Fatal(err.Error())
	}
}
