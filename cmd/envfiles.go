package cmd

import (
	"fmt"
	"errors"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var envfilesCmd = &cobra.Command{
	Use:      "envfiles",
	Short:    "Generate per instance environment json files for deployment",
	Long:     `Generate per instance environment json files for deployment`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pipeline := args[0]

		if viper.GetBool("EnvFilesCreate") {
			generatePipelineEnvfiles(pipeline)
		}
	},
}

func init() {
	createCmd.AddCommand(envfilesCmd)
}

func generatePipelineEnvfiles(pipeline string) {
	var pipeControl PipelineMap

	exitOnError(pipeControl.NewFromFile(pipeControlFile()))

	// load pipeControlFile, error if the requested pipeline is not defined
	pipelineDef, exists := pipeControl[pipeline]
	if !exists { exitOnError(fmt.Errorf("%s not found in %s", pipeline, pipeControlFile())) }

	// start from envDefaultsFile, default file is required even if blank
	baseVars, err := envFileValues(envDefaultsFile())
	exitOnError(err)

	// merge pipeline envfile, if any
	pipelineVars, err := envFileValues(envPipelineFile(pipeline))
	exitOnError(err)
	maps.Copy(baseVars,pipelineVars)

	for _, role := range pipelineDef.Deploy {
		var roleVars = make(map[string]interface{})
		// for each loop, start from baseVars
		maps.Copy(roleVars, baseVars)
		// merge a roles envfile, if any
		maps.Copy(roleVars, anyStaticEnvFileValues(role))
		// write the role envfile to be used by *roleOnly, if any
		writeEnvFileValues(roleVars, role)
		fmt.Printf("writing %s envfile\n", role)
		for _, instance := range pipelineDef.Roles[role].Deploy {
			var instanceVars = make(map[string]interface{})
			maps.Copy(instanceVars, roleVars)
			// merge an instance envfile, if any
			maps.Copy(instanceVars, anyStaticEnvFileValues(instance))
			// merge instance values from pipeControlFile if any
			switch filepath.Ext(viper.GetString("PipeControlFileName"))  {
			case ".json":
				for k, v := range pipelineDef.Roles[role].Instances[instance].(map[string]interface{}) {
					instanceVars[k] = v.(string)
				}
			case ".yaml", ".yml":
				for k, v := range pipelineDef.Roles[role].Instances[instance].(map[interface{}]interface{}) {
					instanceVars[k.(string)] = v
				}
			default:
				exitOnError(errors.New("unsupport envfiles extension"))
			}
			// write the instance envfile
			writeEnvFileValues(instanceVars, instance)
			fmt.Printf("writing %s envfile\n", instance)
		}

	}
}
