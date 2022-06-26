package cmd

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/aboxofsox/optics/pkg/optics"
	"github.com/spf13/cobra"
)

var skip bool

var toml = `name = "{{.Name}}"
scheme = "{{.Scheme}}"
host = "{{.Host}}"
endpoints = [{{range $index, $e := .Endpoints}} "{{$e}}",{{end}}]
outfile = {{.OutFile}}
outdir = "{{.Outdir}}"

[query_params]
# query = ""

[headers]
# Authorization = ""
`

func init() {
	rootCmd.AddCommand(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&skip, "skip", "y", false, "Use an empy config.")
}

var initConfig = &cobra.Command{
	Use:   "init",
	Short: "Initialize optics.toml",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			config    optics.Config
			count     int
			endpoints []string
		)
		p := filepath.Join(".", string(filepath.Separator), "optics.toml")

		if !skip {
			print("> Name: ")
			fmt.Scanln(&config.Name)
			print("> Scheme (http/https): ")
			fmt.Scanln(&config.Scheme)
			print("> Host: ")
			fmt.Scanln(&config.Host)
			print("> Number of endpoints: ")
			fmt.Scanln(&count)
			endpoints = make([]string, count)
			for i := 0; i < count; i++ {
				fmt.Printf("\t> Endpoint %d: ", i+1)
				fmt.Scanln(&endpoints[i])
			}

			print("> Outfile (bool): ")
			fmt.Scanln(&config.OutFile)
			print("> Outdir: ")
			fmt.Scanln(&config.Outdir)
		}

		config.Endpoints = append(config.Endpoints, endpoints...)

		file, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0665)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		tmpl, err := template.New("toml").Parse(toml)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := tmpl.Execute(file, config); err != nil {
			log.Fatal(err.Error())
		}
	},
}
