package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var EnvFilesCreate bool
var EnvFilesPath string
var PipePath string
var PipeControlFileName string
var EnvFilesExt string
var EnvFilesWriteExt string
var PipeOutFile string
var PipeWorkflowName string
var PipeIsApprove bool
var PipeApproveAfterPre bool
var PipeApproveAfterPost bool
var PipePriorJobsRequired bool
var PipeSkipApproval string
var PipeIsPre bool
var PipeIsPost bool
var PipePreRoleOnly bool
var PipePostRoleOnly bool
var PipePreTemplate string
var PipePostTemplate	string
var PipePreJobName string
var PipePostJobName string
var PipeApprovalJobName string
var EnvDefaultsFileName string
var PipeWorkflowHeading string
var PipeApprovalTemplate string
var CircleciConfigFile string

var createCmd = &cobra.Command{
	Use:               "create",
	Short:             "Create circlepipe resources",
	Long:              `Create circlepipe resources

Create the requested circlepipe resource files.

Examples:
  # create the circleci config.yml for the sandbox pipeline
  circlepipe create pipeline sandbox

  # create the instance tfvars files for use in a terraform pipeline
  circlepipe create envfiles

  # create both tfvars and the circleci config.yml
  circlepipe create pipeline sandbox --EnvFilesCreate true
	`,
	ValidArgs: []string{"pipeline", "envfiles"},
	Args: cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		exitOnError(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// flag overrides

  createCmd.PersistentFlags().StringVar(&PipePath, "PipePath", DefaultPipePath, "Path to .circleci pipeline config files")
  if !createCmd.PersistentFlags().Lookup("PipePath").Changed {
    exitOnError(viper.BindPFlag("PipePath", createCmd.PersistentFlags().Lookup("PipePath")))
  }

  rootCmd.PersistentFlags().StringVar(&PipeControlFileName, "PipeControlFileName", DefaultPipeControlFileName, "Name of circlepipe pipeline definition and control file")
  if !rootCmd.PersistentFlags().Lookup("PipeControlFileName").Changed {
    exitOnError(viper.BindPFlag("PipeControlFileName", rootCmd.PersistentFlags().Lookup("PipeControlFileName")))
  }

  rootCmd.PersistentFlags().StringVar(&PipeOutFile, "PipeOutFile", DefaultPipeOutFile, "Name of generated pipeline")
  if !rootCmd.PersistentFlags().Lookup("PipeOutFile").Changed {
    exitOnError(viper.BindPFlag("PipeOutFile", rootCmd.PersistentFlags().Lookup("PipeOutFile")))
  }

  rootCmd.PersistentFlags().StringVar(&PipeWorkflowName, "PipeWorkflowName", DefaultPipeWorkflowName, "Name of workflow in generated pipeline")
  if !rootCmd.PersistentFlags().Lookup("PipeWorkflowName").Changed {
    exitOnError(viper.BindPFlag("PipeWorkflowName", rootCmd.PersistentFlags().Lookup("PipeWorkflowName")))
  }

  createCmd.PersistentFlags().BoolVar(&PipeIsApprove, "PipeIsApprove", DefaultPipeIsApprove, "Include approval step?")
  if !createCmd.PersistentFlags().Lookup("PipeIsApprove").Changed {
    exitOnError(viper.BindPFlag("PipeIsApprove", createCmd.PersistentFlags().Lookup("PipeIsApprove")))
  }

  createCmd.PersistentFlags().BoolVar(&PipeApproveAfterPre, "PipeApproveAfterPre", DefaultPipeApproveAfterPre, "Include approval step after Pre?")
  if !createCmd.PersistentFlags().Lookup("PipeApproveAfterPre").Changed {
    exitOnError(viper.BindPFlag("PipeApproveAfterPre", createCmd.PersistentFlags().Lookup("PipeApproveAfterPre")))
  }

  createCmd.PersistentFlags().BoolVar(&PipeApproveAfterPost, "PipeApproveAfterPost", DefaultPipeApproveAfterPost, "Include approval step after Post?")
  if !createCmd.PersistentFlags().Lookup("PipeApproveAfterPost").Changed {
    exitOnError(viper.BindPFlag("PipeApproveAfterPost", createCmd.PersistentFlags().Lookup("PipeApproveAfterPost")))
  }

  createCmd.PersistentFlags().BoolVar(&PipePriorJobsRequired, "PipePriorJobsRequired", DefaultPipePriorJobsRequired, "Should post jobs require the approval?")
  if !createCmd.PersistentFlags().Lookup("PipePriorJobsRequired").Changed {
    exitOnError(viper.BindPFlag("PipePriorJobsRequired", createCmd.PersistentFlags().Lookup("PipePriorJobsRequired")))
  }

  rootCmd.PersistentFlags().StringVar(&PipeSkipApproval, "PipeSkipApproval", DefaultPipeSkipApproval, "Name of workflow in generated pipeline")
  if !rootCmd.PersistentFlags().Lookup("PipeSkipApproval").Changed {
    exitOnError(viper.BindPFlag("PipeSkipApproval", rootCmd.PersistentFlags().Lookup("PipeSkipApproval")))
  }

  createCmd.PersistentFlags().BoolVar(&PipeIsPre, "PipeIsPre", DefaultPipeIsPre, "Create pre-approval jobs?")
  if !createCmd.PersistentFlags().Lookup("PipeIsPre").Changed {
    exitOnError(viper.BindPFlag("PipeIsPre", createCmd.PersistentFlags().Lookup("PipeIsPre")))
  }

  createCmd.PersistentFlags().BoolVar(&PipeIsPost, "PipeIsPost", DefaultPipeIsPost, "Create post-approval jobs?")
  if !createCmd.PersistentFlags().Lookup("PipeIsPost").Changed {
    exitOnError(viper.BindPFlag("PipeIsPost", createCmd.PersistentFlags().Lookup("PipeIsPost")))
  }

  createCmd.PersistentFlags().BoolVar(&PipePreRoleOnly, "PipePreRoleOnly", DefaultPipePreRoleOnly, "Create create pre-approval jobs only at the role level?")
  if !createCmd.PersistentFlags().Lookup("PipePreRoleOnly").Changed {
    exitOnError(viper.BindPFlag("PipePreRoleOnly", createCmd.PersistentFlags().Lookup("PipePreRoleOnly")))
  }

  createCmd.PersistentFlags().BoolVar(&PipePostRoleOnly, "PipePostRoleOnly", DefaultPipePostRoleOnly, "Create create post-approval jobs only at the role level?")
  if !createCmd.PersistentFlags().Lookup("PipePostRoleOnly").Changed {
    exitOnError(viper.BindPFlag("PipePostRoleOnly", createCmd.PersistentFlags().Lookup("PipePostRoleOnly")))
  }

  createCmd.PersistentFlags().StringVar(&PipePreTemplate, "PipePreTemplate", DefaultPipePreTemplate, "Name of the pre-approval template file")
  if !createCmd.PersistentFlags().Lookup("PipePreTemplate").Changed {
    exitOnError(viper.BindPFlag("PipePreTemplate", createCmd.PersistentFlags().Lookup("PipePreTemplate")))
  }

  createCmd.PersistentFlags().StringVar(&PipePostTemplate, "PipePostTemplate", DefaultPipePostTemplate, "Name of the post-approval template file")
  if !createCmd.PersistentFlags().Lookup("PipePostTemplate").Changed {
    exitOnError(viper.BindPFlag("PipePostTemplate", createCmd.PersistentFlags().Lookup("PipePostTemplate")))
  }

  createCmd.PersistentFlags().StringVar(&PipePreJobName, "PipePreJobName", DefaultPipePreJobName, "Name for the pre-approval jobs, must include {{.instance}}")
  if !createCmd.PersistentFlags().Lookup("PipePreJobName").Changed {
    exitOnError(viper.BindPFlag("PipePreJobName", createCmd.PersistentFlags().Lookup("PipePreJobName")))
  }

  createCmd.PersistentFlags().StringVar(&PipePostJobName, "PipePostJobName", DefaultPipePostJobName, "Name for the post-approval jobs, must include {{.instance}}")
  if !createCmd.PersistentFlags().Lookup("PipePostJobName").Changed {
    exitOnError(viper.BindPFlag("PipePostJobName", createCmd.PersistentFlags().Lookup("PipePostJobName")))
  }

  createCmd.PersistentFlags().StringVar(&PipeApprovalJobName, "PipeApprovalJobName", DefaultPipeApprovalJobName, "Name for the approval jobs, must include {{.jobname}}")
  if !createCmd.PersistentFlags().Lookup("PipeApprovalJobName").Changed {
    exitOnError(viper.BindPFlag("PipeApprovalJobName", createCmd.PersistentFlags().Lookup("PipeApprovalJobName")))
  }

  createCmd.PersistentFlags().StringVar(&EnvDefaultsFileName, "EnvDefaultsFileName", DefaultEnvDefaultsFileName, "Name of the envfile containing the pipeline-wide default values")
  if !createCmd.PersistentFlags().Lookup("EnvDefaultsFileName").Changed {
    exitOnError(viper.BindPFlag("EnvDefaultsFileName", createCmd.PersistentFlags().Lookup("EnvDefaultsFileName")))
  }

  createCmd.PersistentFlags().StringVar(&PipeWorkflowHeading, "PipeWorkflowHeading", DefaultPipeWorkflowHeading, "(advanced) name of the workflow setup template")
  if !createCmd.PersistentFlags().Lookup("PipeWorkflowHeading").Changed {
    exitOnError(viper.BindPFlag("PipeWorkflowHeading", createCmd.PersistentFlags().Lookup("PipeWorkflowHeading")))
  }

  createCmd.PersistentFlags().StringVar(&PipeApprovalTemplate, "PipeApprovalTemplate", DefaultPipeApprovalTemplate, "(advanced) name of the approval step template")
  if !createCmd.PersistentFlags().Lookup("PipeApprovalTemplate").Changed {
    exitOnError(viper.BindPFlag("PipeApprovalTemplate", createCmd.PersistentFlags().Lookup("PipeApprovalTemplate")))
  }

  createCmd.PersistentFlags().StringVar(&CircleciConfigFile, "CircleciConfigFile", DefaultCircleciConfigFile, "(advanced) alternate name for the .circleci config.yml")
  if !createCmd.PersistentFlags().Lookup("CircleciConfigFile").Changed {
    exitOnError(viper.BindPFlag("CircleciConfigFile", createCmd.PersistentFlags().Lookup("CircleciConfigFile")))
  }

  createCmd.PersistentFlags().BoolVar(&EnvFilesCreate, "EnvFilesCreate", DefaultEnvFilesCreate, "Generate instance level environment json files")
  if !createCmd.PersistentFlags().Lookup("EnvFilesCreate").Changed {
    exitOnError(viper.BindPFlag("EnvFilesCreate", createCmd.PersistentFlags().Lookup("EnvFilesCreate")))
  }

  createCmd.PersistentFlags().StringVar(&EnvFilesPath, "EnvFilesPath", DefaultEnvFilesPath, "Path to environment json files")
  if !createCmd.PersistentFlags().Lookup("EnvFilesPath").Changed {
    exitOnError(viper.BindPFlag("EnvFilesPath", createCmd.PersistentFlags().Lookup("EnvFilesPath")))
  }

  createCmd.PersistentFlags().StringVar(&EnvFilesExt, "EnvFilesExt", DefaultEnvFilesExt, "Environment file format")
  if !createCmd.PersistentFlags().Lookup("EnvFilesExt").Changed {
    exitOnError(viper.BindPFlag("EnvFilesExt", createCmd.PersistentFlags().Lookup("EnvFilesExt")))
  }

  createCmd.PersistentFlags().StringVar(&EnvFilesWriteExt, "EnvFilesWriteExt", DefaultEnvFilesWriteExt, "Generated environment file format")
  if !createCmd.PersistentFlags().Lookup("EnvFilesWriteExt").Changed {
    exitOnError(viper.BindPFlag("EnvFilesWriteExt", createCmd.PersistentFlags().Lookup("EnvFilesWriteExt")))
  }

}
