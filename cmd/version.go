package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:               "version",
	Short:             "Show circlepipe version information",
	Long:              `Show circlepipe version information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("circlepipe " + Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
