package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"errors"
	"text/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/maps"
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Generate config.yml for specified pipeline",
	Long:  `Generate config.yml for specified pipeline`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var jobsToBeApproved, requiredJobs, lastJob string
		var pipeControl PipelineMap

		pipeline := args[0]
		exitOnError(pipeControl.NewFromFile(pipeControlFile()))

		// load pipeControlFile, error if the requested pipeline is not defined
		pipelineDef, exists := pipeControl[pipeline]
		if !exists { exitOnError(fmt.Errorf("%s not found in %s", pipeline, pipeControlFile())) }

		// get the pipeline filter, error if no filter is defined
		filter := pipelineDef.Filter
		if filter == "" { exitOnError(fmt.Errorf("filter not found in %s pipeline", pipeline)) }

		// also generate envfiles if EnvFilesCreate is true
		if viper.GetBool("EnvFilesCreate") {
			fmt.Println("EnvFileCreate flag set.")
			generatePipelineEnvfiles(pipeline)
		}

		fmt.Printf("generating pipeline: %s\n\n", pipeline)
		fmt.Printf("  filter: %s\n", pipelineDef.Filter)

		// copy everything but the jobs and workflows from CircleciConfigFile into PipeOutFile
		setupGeneratedConfigOutfile()

		// prep the templates
		pre, approve, post := loadTemplates()

		// open PipeOutFile for append
		outfile, err := os.OpenFile(generatedConfigFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		exitOnError(err)
		defer outfile.Close()

		// set the starting point for requires/approval related config that changes throughout the iteration
		jobsToBeApproved = "requires:"
		requiredJobs = "requires:"
		lastJob = ""

		for _, role := range pipelineDef.Deploy {
			fmt.Printf("  role: %s\n", role)

			// generate jobs that come before an optional approval step
			if viper.GetBool("PipeIsPre") {

				// optionally, the pipelines can be built around steps that occur only per role in PipeControl
				if viper.GetBool("PipePreRoleOnly") {
					jobsToBeApproved += jobRequiresLine("pre", role)
					roleVars := asssembleVars(filter, role, "", jobNameTemplate("pre", role), lastJob, requiredJobs, "", "", anyEnvFileValues(role))
					exitOnError(pre.Execute(outfile, roleVars))

				} else {
					// default, pipelines built around steps that occur per instance in PipeControl
					for _, instance := range pipelineDef.Roles[role].Deploy {
						fmt.Printf("    [pre] %s\n", instance)
						jobsToBeApproved += jobRequiresLine("pre", instance)
						instanceVars := asssembleVars(filter, role, instance, jobNameTemplate("pre", role), lastJob, requiredJobs, "", "", anyEnvFileValues(instance))
						exitOnError(pre.Execute(outfile, instanceVars))
					}
				}
			}
			//
			lastJob = "requires:" + jobRequiresLine("pre", role)
			// generate an optional approval step
			requiredJobs = "requires:"
			if viper.GetBool("PipeIsApprove") {
				fmt.Printf("      [approve] %s\n", role)
				requiredJobs += jobRequiresLine("approve", role)
				approvalVars := asssembleVars(filter, role, "", jobNameTemplate("approve", role), lastJob, requiredJobs, "", jobsToBeApproved, make(map[string]interface{}))
				exitOnError(approve.Execute(outfile, approvalVars))
			}
			jobsToBeApproved = "requires:"

			// generate jobs that come after an optional approval step
			if viper.GetBool("PipeIsPost") {

				// optionally, the pipelines can be built around steps that occur only per role in PipeControl
				if viper.GetBool("PipePreRoleOnly") {
					jobsToBeApproved += jobRequiresLine("post", role)
					roleVars := asssembleVars(filter, role, "", jobNameTemplate("post", role), lastJob, requiredJobs, jobNameTemplate("approve", role), "", anyEnvFileValues(role))
					exitOnError(post.Execute(outfile, roleVars))

				} else {
					// default, pipelines built around steps that occur per instance in PipeControl
					for _, instance := range pipelineDef.Roles[role].Deploy {
						fmt.Printf("    [post] %s\n", instance)

						jobsToBeApproved += jobRequiresLine("post", instance)
						instanceVars := asssembleVars(filter, role, instance, jobNameTemplate("post", role), lastJob, requiredJobs, jobNameTemplate("approve", role), "", anyEnvFileValues(instance))
						exitOnError(post.Execute(outfile, instanceVars))
					}
				}
			}
		}
	},
}

func init() {
	createCmd.AddCommand(pipelineCmd)
}

func asssembleVars(filter interface{}, role, instance, jobName, lastJob, requiredJobs, approvalJobName, jobsToBeApproved string, fromEnv map[string]interface{}) map[string]interface{} {
	jsonVars := make(map[string]interface{})
	jsonVars["filter"] = filter
	jsonVars["role"] = role
	jsonVars["instance"] = instance
	jsonVars["envfilespath"] = viper.GetString("EnvFilesPath")
	jsonVars["jobname"] = jobName
	jsonVars["approvaljobname"] = approvalJobName
	jsonVars["jobstobeapproved"] = jobsToBeApproved
	if lastJob != "" {
		jsonVars["lastjob"] = lastJob
	} else {
		jsonVars["lastjob"] = ""
	}
	maps.Copy(jsonVars, fromEnv)
	if viper.GetBool("PipePriorJobsRequired") && requiredJobs != "requires:" {
		jsonVars["priorapprovalrequired"] = requiredJobs
	} else {
		jsonVars["priorapprovalrequired"] = ""
	}
	return jsonVars
}

func loadTemplates() (*template.Template, *template.Template, *template.Template) {
	// load the templates into text/template format for use
	preContents := preTemplateFile()
	approvalContents := approvalTemplate()
	postContents := postTemplateFile()

	pre, err := template.New("preTemplate").Parse(preContents)
	exitOnError(err)

	approval, err := template.New("approvalTemplate").Parse(approvalContents)
	exitOnError(err)

	post, err := template.New("postTemplate").Parse(postContents)
	exitOnError(err)

	return pre, approval, post
}

func approvalTemplate() string {
	if viper.GetString("PipeApprovalTemplate") == "default" {
		return DefaultApproval
	} else {
		return textFromFile(viper.GetString("PipePath") + "/" + viper.GetString("PipeApprovalTemplate"))
	}
}

func jobRequiresLine(jobName string, jobType string) string {
	return fmt.Sprintf("\n            - %s", jobNameTemplate(jobName, jobType))
}

func jobNameTemplate(jobName string, jobType string) string {
	var result string
	if jobName == "pre" {
		result = fmt.Sprintf(viper.GetString("PipePreJobName"), jobType)
	} else if jobName == "post" {
		result = fmt.Sprintf(viper.GetString("PipePostJobName"), jobType)
	} else if jobName == "approve" {
		result = fmt.Sprintf(viper.GetString("PipeApprovalJobName"), jobType)
	} else {
		exitOnError(errors.New("only pre or post jobNames is supported"))
	}
	return result
}

func setupGeneratedConfigOutfile() {
	// create starting point for PipeOutFile
	lines := generateConfigLines()
	outFile, err := os.Create(generatedConfigFile())
	exitOnError(err)
	defer outFile.Close()

	for _, line := range lines {
		_, err := outFile.WriteString(line + "\n")
		exitOnError(err)
	}
	_, err = outFile.WriteString(fmt.Sprintf(workflowHeading(), viper.GetString("PipeWorkflowName")))
	exitOnError(err)
}

func generateConfigLines() []string {
	// returns the necessary portion of CircleciConfigFile
	lines := []string{}
	readUntilJobsOrWorkflows(circleciConfigFile(), func(line string) {
		lines = append(lines, line)
	})
	return lines
}

func readUntilJobsOrWorkflows(filePath string, callback func(string)) {
	// read from CircleciConfigFile until jobs: or workflows:
	// also remove the setup: directive needed by calling pipeline
	file, err := os.Open(filePath)
	exitOnError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "setup:") {
			continue
		}
		if strings.HasPrefix(line, "workflows:") {
			break
		}
		callback(line)
	}
	exitOnError(scanner.Err())
}

func workflowHeading() string {
	// return the default workflow heading
	if viper.GetString("PipeWorkflowHeading") == "default" {
		return DefaultWorkflowHeading
	} else {
		return textFromFile(viper.GetString("PipePath") + "/" + viper.GetString("PipeWorkflowHeading"))
	}
}
