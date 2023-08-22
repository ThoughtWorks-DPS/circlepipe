package cmd

import (
	"github.com/spf13/viper"
)

func pipeControlFile() string {
	return viper.GetString("EnvFilesPath") + "/" + viper.GetString("PipeControlFileName")
}

func circleciConfigFile() string {
	return viper.GetString("PipePath") + "/" + viper.GetString("CircleciConfigFile")
}

func generatedConfigFile() string {
  return viper.GetString("PipePath") + "/" + viper.GetString("PipeOutFile")
}

func preTemplateFile() string {
	return textFromFile(viper.GetString("PipePath") + "/" + viper.GetString("PipePreTemplate"))
}

func postTemplateFile() string {
	return textFromFile(viper.GetString("PipePath") + "/" + viper.GetString("PipePostTemplate"))
}

func envDefaultsFile() string {
	return viper.GetString("EnvFilesPath") + "/" + viper.GetString("EnvDefaultsFileName") + viper.GetString("EnvFilesExt")
}

func envPipelineFile(pipeline string) string {
	return viper.GetString("EnvFilesPath") + "/" +pipeline + viper.GetString("EnvFilesExt")
}

func envStaticFile(fn string) string {
	return viper.GetString("EnvFilesPath") + "/" + fn + viper.GetString("EnvFilesExt")
}

func envFile(fn string) string {
	return viper.GetString("EnvFilesPath") + "/" + fn + viper.GetString("EnvFilesWriteExt")
}
