package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

type (
	rootGenerateCmdOptions struct {
	}
)

var rootGenerateOptions rootGenerateCmdOptions

var rootGenerateCmd = &cobra.Command{
	Use: "generate",

	RunE: func(cmd *cobra.Command, args []string) error {
		return execrootGenerate(rootGenerateOptions)
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		rootCmd.PreRun(cmd, args)

	},
}

func init() {
	rootCmd.AddCommand(rootGenerateCmd)
}

func execrootGenerate(options rootGenerateCmdOptions) error {
	fmt.Println("Running cmd: generate")
	fmt.Printf("rootGenerateCmdOptions: %v\n", options)
	return nil
}
