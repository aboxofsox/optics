package optics

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

func Test_cPath(t *testing.T) {
	p := cPath(".", "/", "test")
	if p != "./test" {
		t.Errorf("got: %s, want %s\n", p, "./test")
	}
}

func Test_read(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.Write([]byte("hello"))
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{args: args{
			r: buf,
		}},
	}
	for _, tt := range tests {
		tt.want = []byte("hello")
		t.Run(tt.name, func(t *testing.T) {
			got, err := read(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_opTime(t *testing.T) {
	type args struct {
		start time.Time
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{
			start: time.Now(),
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := opTime(tt.args.start); got != tt.want {
				t.Errorf("opTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	// replace test cases with your own, or use this one.
	c := &Config{
		Name:        "pokeapi",
		Scheme:      "https",
		Host:        "pokeapi.co",
		Endpoints:   []string{"api/v2/pokemon/pikachu", "api/v2/pokemon/mew"},
		OutFile:     true,
		Outdir:      "res",
		QueryParams: map[string]string{},
		Headers:     map[string]string{},
	}

	s := `name = "pokeapi"
scheme = "https"
host = "pokeapi.co"
endpoints = [ "api/v2/pokemon/pikachu", "api/v2/pokemon/mew",]
outfile = true
outdir = "res"

[query_params]
# query = ""

[headers]
# Authorization = ""`

	d, _ := decode(s)

	if !reflect.DeepEqual(c, d) {
		t.Error("struct does not match test struct.")
	}
}
