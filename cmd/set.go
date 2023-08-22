package cmd

import (
	//"fmt"
	"os"
	"errors"
	"regexp"
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:               "set",
	Short:             "Set a specific circlepipe configuration key setting",
	Long:              `Set a specific circlepipe configuration key setting

Example:
  # set generate pipeline filename equal to deploy.yml
  circlepipe config set PipeOutFile deploy.yml
	`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if viper.IsSet(args[0]) {
			exitOnError(validateSetting(args))
			viper.Set(args[0], args[1])
			exitOnError(viper.WriteConfig())
		} else {
			exitOnError(cmd.Help())
			exitOnError(errors.New("error: configuration key not a supported customization"))
		}
	},
}

func init() {
	configCmd.AddCommand(setCmd)
}

func validateSetting(args []string) error {
	switch strings.ToLower(args[0]) {

	// must be valid workflow or job name
	case "pipeprejobname",
			 "pipeworkflowname",
			 "pipepostjobname",
			 "pipeapprovaljobname":
		if len(args[1]) == 0 || len(args[1]) > 128 || ! regexp.MustCompile(`^[a-zA-Z -]+$`).MatchString(args[1]) {
			return errors.New("circlepipe set error: workflow name too long or contains invalid characters")
		}

	// must be boolean
	case "pipeisapprove",
		   "pipepriorjobsrequired",
			 "pipeispre",
			 "pipeispost",
			 "pipepreroleonly",
			 "pipepostroleonly",
			 "envfilescreate",
			 "pipecreateapprovalstep",
			 "piperoleonly":
		if args[1] != "true" && args[1] != "false" {
			 return errors.New("circlepipe set error: boolean key value setting must be true or false")
		}

	// must be valid, existing path and/or file
	case "pipepretemplate",
			 "pipeposttemplate",
			 "envfilespath",
			 "envdefaultsfilename",
			 "pipepath",
			 "pipecontrolfilename",
			 "pipeapprovaltemplate":
		_, err := os.Stat(args[1])
		return err

	// must be a valid path and/or file name (will be created so doesn't check if alredy exists)
	case "pipeoutfile":
		if len(args[1]) == 0 || len(args[1]) > 255 || ! regexp.MustCompile(`^[^\\/:\*\?"<>\|]+$`).MatchString(args[1]) {
			return errors.New("circlepipe set error: invalid outfile name")
		}

	// extensions must start with '.'
	case "envfilesext",
			 "envfileswriteext":
		if len(args[1]) == 0 || len(args[1]) > 128 || ! regexp.MustCompile(`^[.][^\\/:\*\?"<>\|]*$`).MatchString(args[1]) {
			return errors.New("circlepipe set error: extensions must start with '.'")
		}
	default:
		return errors.New("attempted to set unknown setting")
	}
	return nil
}
