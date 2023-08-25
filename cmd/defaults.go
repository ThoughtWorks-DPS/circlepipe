package cmd

// circlepipe default settings
const (

	// Default name of file containing the pipeline definition
	DefaultPipeControlFileName	  = "generate.yaml"

	// Default name of the generated pipeline config.yml
	DefaultPipeOutFile  				  = "generated_config.yml"

	// default name for the generated workflow
	DefaultPipeWorkflowName       = "continuation-generated-workflow"

	// Include an approval step between pro and post jobs
	DefaultPipeIsApprove				  = true
	// Should jobs require the prior approval step, usually set to match PipeIsApprove
	DefaultPipePriorJobsRequired  = true
	// Skip approval step for this role
	DefaultPipeSkipApproval				= ""

	// Include Pre or Post jobs
	DefaultPipeIsPre              = true
	DefaultPipeIsPost             = true

	// Should there be Pre or Post jobs for each role rather than each instance
	DefaultPipePreRoleOnly        = false
	DefaultPipePostRoleOnly       = false

	// Default names for pre, approve, and post job template files
	DefaultPipePreTemplate				= "pre-approve.yml"
	DefaultPipePostTemplate				= "post-approve.yml"

	// Default names for pre, approve, and post jobs
	DefaultPipePreJobName         = "plan %s change"
	DefaultPipePostJobName        = "apply %s change"
	DefaultPipeApprovalJobName    = "approve %s changes"

	// Generate environment setting json files
	DefaultEnvFilesCreate         = true

	// Default folder for environment setting files
	DefaultEnvFilesPath           = "environments"

	// Default extension for environment setting files
	// the files will be json format regardless of this setting
	DefaultEnvFilesExt						= ".json"

	// Default name
  DefaultEnvDefaultsFileName    = "default"
	DefaultEnvFilesWriteExt       = ".tfvars.json"

	// Advanced configuration options, not typically used
	// Approval template is defined within the cli by default but can override by setting to file reference
	DefaultPipeWorkflowHeading    = "default"
	DefaultPipeApprovalTemplate   = "default"
	DefaultPipePath         		  = ".circleci"
	DefaultCircleciConfigFile     = "config.yml"

	// settings for cli configuration file
	ConfigEnvDefault              = "CIRCLEPIPE"
	ConfigFileDefaultName         = ".circlepipe"
	ConfigFileDefaultType         = "yaml"
	ConfigFileDefaultLocation     = "."
	ConfigFileDefaultLocationMsg  = "config file (default is ./.circlepipe.yaml)"
)

const (

	DefaultWorkflowHeading = `
workflows:
  version: 2

  %s:
    jobs:
`
	DefaultApproval = `
      - {{.jobname}}:
          type: approval
          {{.jobstobeapproved}}
          filters: {{.filter}}

`
)
