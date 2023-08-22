package cmd

import (
	"fmt"
	"encoding/json"
	"sort"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewCmd = &cobra.Command{
	Use:    "view",
	Short:  "View circlepipe configuration",
	Long:   `View the circlepipe configuration of all customizable settings.`,
	Run: func(cmd *cobra.Command, args []string) {

		configurations, err := formattedConfigurationList()
		exitOnError(err)
		fmt.Println(configurations)
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}

func formattedConfigurationList() (string, error) {
	keysJson := getConfiguration()
	formattedConfig, err := json.MarshalIndent(keysJson, "", "  ")
	return string(formattedConfig), err
}

func getConfiguration() map[string]interface{} {
	var keysJson = make(map[string]interface{})

	keys := viper.AllKeys()
	sort.Strings(keys)
	for _, key := range keys {
		keysJson[key] = viper.Get(key)
	}
	fmt.Println(len(keysJson))
	return keysJson
}
