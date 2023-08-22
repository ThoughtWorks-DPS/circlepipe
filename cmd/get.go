package cmd

import (
	"fmt"
	"errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCmd = &cobra.Command{
	Use:               "get",
	Short:             "Get a specific circlepipe configuration key setting",
	Long:              `Get a specific circlepipe configuration key setting

Example:
  $ circlepipe config get PipePath
  .circleci
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet(args[0]) {
			fmt.Printf("%v", viper.Get(args[0]))
		} else {
			exitOnError(cmd.Help())
			exitOnError(errors.New("error: configuration key not found"))
		}
	},
}

func init() {
	configCmd.AddCommand(getCmd)
}
