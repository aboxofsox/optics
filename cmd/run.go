package cmd

import (
	"fmt"

	"github.com/aboxofsox/optics/pkg/optics"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runOptics)
}

var runOptics = &cobra.Command{
	Use:   "run",
	Short: "Run optics",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running optics")
		optics := optics.New()
		optics.Init()
	},
}
