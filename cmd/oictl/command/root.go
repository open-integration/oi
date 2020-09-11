package command

import (
	"github.com/spf13/viper"
	"os"

	"fmt"
	"github.com/spf13/cobra"
)

var cnf *viper.Viper = viper.New()

type (
	rootCmdOptions struct {
		verbose bool
	}
)

var rootOptions rootCmdOptions

var rootCmd = &cobra.Command{
	Use:     "oictl",
	Version: "0.1.0",
	PreRun: func(cmd *cobra.Command, args []string) {

		cnf.Set("verbose", rootOptions.verbose)

	},
}

// Execute - execute the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().BoolVar(&rootOptions.verbose, "verbose", cnf.GetBool("verbose"), "")
}
