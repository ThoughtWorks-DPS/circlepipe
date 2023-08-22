package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:               "init",
	Short:             "Create .circlepipe.yaml configuration files from current settings",
	Long:              `Create .circlepipe.yaml configuration files from current settings`,
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(viper.WriteConfig())
		fmt.Println(viper.GetString("PipeWorkflowHeading"))
		fmt.Println("current configuration written to config file")
	},
}

func init() {
	configCmd.AddCommand(initCmd)
}
