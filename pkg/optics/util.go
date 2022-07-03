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

// Handle HTTP status codes
var StatusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",

	200: "Ok",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Counter",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",

	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "USe Proxy",
	306: "Unused",
	307: "Temporary Redirect",
	308: "Permanent Redirect",

	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Found",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Miscreated Request",
	422: "Unprocessable Entiy",
	423: "Locked",
	424: "Failed Dependency",
	425: "Too Early",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",

	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

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

func opTime(start time.Time) time.Duration  { return time.Since(start) }

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
