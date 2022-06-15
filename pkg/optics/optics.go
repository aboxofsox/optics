package optics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
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

type Controller struct {
	Config       *Config
	Buffer       *bytes.Buffer
	Client       *http.Client
	Url          *url.URL
	HttpResponse *HttpResponse
	Data         []byte
}

func New() *Controller {
	w := &strings.Builder{}
	var (
		envMap map[string]string
		config *Config
	)

	if envExists() {
		envMap, _ = godotenv.Read()
	}

	configPath := filepath.Join(".", string(filepath.Separator), "optics.toml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal(err.Error())
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	f, _ := os.Stat(configPath)
	if f.Size() <= 0 {
		log.Fatal("empty config file")
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	tmpl, err := template.New("toml").Parse(string(data))
	if err != nil {
		log.Fatal(err.Error())
	}

	tmpl.Execute(w, envMap)

	if _, err := toml.Decode(w.String(), &config); err != nil {
		log.Fatal(err.Error())
	}

	return &Controller{
		Config: config,
		Buffer: &bytes.Buffer{},
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		HttpResponse: &HttpResponse{},
	}

}

func (ctrl *Controller) Init() {
	var wg sync.WaitGroup
	ctrl.Url = &url.URL{
		Scheme: ctrl.Config.Scheme,
		Host:   ctrl.Config.Host,
	}

	if _, err := os.Stat(ctrl.Config.Outdir); os.IsNotExist(err) {
		os.Mkdir(ctrl.Config.Outdir, 0665)
	}

	if len(ctrl.Config.Endpoints) == 0 {
		log.Fatal("no endpoints to test")
	}

	for i := range ctrl.Config.Endpoints {
		wg.Add(1)
		// buf := &bytes.Buffer{}

		ctrl.Url.Path = ctrl.Config.Endpoints[i]

		kv := ctrl.Url.Query()

		if len(ctrl.Config.QueryParams) != 0 {
			for k, v := range ctrl.Config.QueryParams {
				kv.Add(k, v)
			}

			ctrl.Url.RawQuery = kv.Encode()
		}

		go ctrl.Get(ctrl.Url.String(), wg.Done)
		wg.Wait()
	}
}

func (ctrl *Controller) Get(url string, done func()) {
	defer done()

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

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	since := time.Since(start)

	ctrl.Json(data)
	ctrl.Log(*res, since)
	ctrl.Buffer.Reset()

}

func (ctrl *Controller) Json(data []byte) (int, error) {
	buf := new(bytes.Buffer)
	if len(data) == 0 {
		return 0, fmt.Errorf("no data to write: %d", len(data))
	}
	if err := json.Indent(buf, data, "", " "); err != nil {
		return 0, err
	}

	tmp, err := ioutil.TempFile(ctrl.Config.Outdir, fmt.Sprintf("%s-*.json", epName(*ctrl.Url)))
	if err != nil {
		return 0, err
	}
	defer tmp.Close()

	if _, err := tmp.Write(buf.Bytes()); err != nil {
		return 0, err
	}

	return len(data), nil

}

func (ctrl *Controller) Log(res http.Response, duration time.Duration) {
	headers := map[string][]string{}
	for k, v := range res.Request.Header {
		headers[k] = append(headers[k], v...)
	}

	ctrl.HttpResponse.StatusCode = res.StatusCode
	ctrl.HttpResponse.Method = res.Request.Method
	ctrl.HttpResponse.Time = duration.Seconds()
	ctrl.HttpResponse.Headers = headers

	lfp := filepath.Join(".", string(filepath.Separator), ctrl.Config.Outdir, "optics.log")

	tm := time.Now()
	timestamp := fmt.Sprintf("%d-%02d-%02d %02d:%02d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute())

	file, err := os.OpenFile(lfp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
		duration.Seconds(),
	)); err != nil {
		log.Fatal(err.Error())
	}
}
