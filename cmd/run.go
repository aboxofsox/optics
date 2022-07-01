package cmd

import (
	"fmt"

	"github.com/aboxofsox/optics/pkg/optics"
	"github.com/spf13/cobra"
)

var (
	useProxy bool
)

func init() {
	rootCmd.AddCommand(runOptics)
	rootCmd.PersistentFlags().BoolVarP(&useProxy, "use-proxy", "p", false, "Proxy endpoint tests.")
}

var runOptics = &cobra.Command{
	Use:   "run",
	Short: "Run optics",
	Run: func(cmd *cobra.Command, args []string) {
		options := &optics.Options{
			UseProxy: useProxy,
		}
		fmt.Println("running optics")
		optics := optics.New()
		optics.Init(options)
	},
}
