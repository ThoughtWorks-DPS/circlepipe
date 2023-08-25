package cmd

import (
  "os"
  "log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Version = "snapshot"

var rootCmd = &cobra.Command{
  Use:   "circlepipe",
  Short: "circlepipe is a circleci pipeline generator",
  Long:  `A fast and flexible pipeline generator designed for use with the circleci cotniuation orb.
Complete documentation is available at https://github.com/ThoughtWorks-DPS/circlepipe`,
  Run: func(cmd *cobra.Command, args []string) {
		exitOnError(cmd.Help())
  },
}

func Execute() {
  exitOnError(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// you may specify the config file and location. Viper supports the following file types based on extension:
	// JSON, TOML, YAML, HCL, INI, envfile and Java Properties files
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", ConfigFileDefaultLocationMsg)
}

// initConfig sets the config values based on the following order of precedent:
// Default values from constant.go
// override with Config file definitions if present
// override with ENV variables if present
func initConfig() {
  viper.SetDefault("CircleciConfigFile", DefaultCircleciConfigFile)
	viper.SetDefault("PipeControlFileName", DefaultPipeControlFileName)
  viper.SetDefault("PipeOutFile", DefaultPipeOutFile)
  viper.SetDefault("PipePath", DefaultPipePath)
  viper.SetDefault("PipeWorkflowName", DefaultPipeWorkflowName)
  viper.SetDefault("PipeIsApprove", DefaultPipeIsApprove)
  viper.SetDefault("PipePriorJobsRequired", DefaultPipePriorJobsRequired)
  viper.SetDefault("PipeSkipApproval", DefaultPipeSkipApproval)
  viper.SetDefault("PipePreRoleOnly", DefaultPipePreRoleOnly)
  viper.SetDefault("PipePostRoleOnly", DefaultPipePostRoleOnly)
  viper.SetDefault("PipeIsPre", DefaultPipeIsPre)
  viper.SetDefault("PipeIsPost", DefaultPipeIsPost)
  viper.SetDefault("PipePreTemplate", DefaultPipePreTemplate)
  viper.SetDefault("PipePostTemplate", DefaultPipePostTemplate)
  viper.SetDefault("PipeWorkflowHeading", DefaultPipeWorkflowHeading)
  viper.SetDefault("PipeApprovalTemplate", DefaultPipeApprovalTemplate)
  viper.SetDefault("PipePreJobName", DefaultPipePreJobName)
  viper.SetDefault("PipePostJobName", DefaultPipePostJobName)
  viper.SetDefault("PipeApprovalJobName", DefaultPipeApprovalJobName)
  viper.SetDefault("EnvFilesCreate", DefaultEnvFilesCreate)
  viper.SetDefault("EnvFilesPath", DefaultEnvFilesPath)
  viper.SetDefault("EnvFilesExt", DefaultEnvFilesExt)
  viper.SetDefault("EnvDefaultsFileName", DefaultEnvDefaultsFileName)
  viper.SetDefault("EnvFilesWriteExt", DefaultEnvFilesWriteExt)

	viper.SetEnvPrefix(ConfigEnvDefault)
	viper.AutomaticEnv()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(ConfigFileDefaultLocation)
		viper.SetConfigName(ConfigFileDefaultName)
	}

	// If a config file is found, read it in, else write a blank.
	if err := viper.ReadInConfig(); err != nil {
    if cfgFile == "" {
      cfgFile = ConfigFileDefaultLocation + "/" + ConfigFileDefaultName + "." + ConfigFileDefaultType
    }
		emptyFile, err := os.Create(cfgFile)
		exitOnError(err)
		emptyFile.Close()
	}
}

func exitOnError(err error) bool {
	if err != nil {
		log.Fatal(err)
	}
	return true
}
