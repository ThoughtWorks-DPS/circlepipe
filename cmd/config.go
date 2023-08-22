package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:               "config",
	Short:             "Modify circlepipe configuration file",
	Long:              `View and make changes to the circlepipe configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
