<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## 6 Customization options

Circlepipe supports a number of customization options. To override the defaults to use custom values include the following flags.  

Alternatively, you may also define a `.circlepipe.yaml` configuration file. See the instructions for the using the config file below. Yaml or Json file format is supported for all files.  

#### 6.1 `--Flags`

###### 6.1.1 `--PipeControlFileName`

**default** _generate.yaml_

Name of the file containing the pipeline control definition.  

#### 6.1.2 `--PipeOutFile`

**default** _generated_config.yml_

Name of the resulting, generated circleci config.yml.

#### 6.1.3 `--PipeWorkflowName`

**default** _"continuation-generated-workflow"_

Name of the workflow in the generated config.yml

#### 6.1.4 `--PipeIsPre`

**default** _true_

Should be a job from the pre-approve template for each role or instance in the control file.  

#### 6.1.4 `--PipeIsApprove`

**default** _true_

Should be an Approval step after the pre- jobs for each role in the control file.  

#### 6.1.4 `--PipeIsPost`

**default** _true_

Should be a job from the post-approve template for each role or instance in the control file.  

#### 6.1.3 `--PipePriorJobsRequired`

**default** _true_

Jobs after approval steps should require an Approval. Where there will be an approval step, this should alwyas be true.

#### 6.1.4 `--PipePreRoleOnly`

**default** _false_

The "pre" jobs created should be at the role-level only. Normally pre- jobs are created for each instance in a role. But for deployment pipelines that do not include any instance specific action you can set this flag to generate an instance of the pre-approve template once per role. See **other pipeline configuration examples** below.

#### 6.1.4 `--PipePostRoleOnly`

**default** _false_

Similar to the previsou flag, there may be a situation in which a post-approve job is required but only at the role-level. Set this flag to get the result.  

#### 6.1.4 `--PipePreTemplate`

**default** _pre-approve.yml_  

Name of the file containing the pre-approve.yml template.


#### 6.1.4 `--PipePostTemplate`

**default** _post-approve.yml_  

Name of the file containing the post-approve.yml template.

#### 6.1.4 `--PipePreJobName`

**default** _"plan %s change"_  

The configuration of the job name for each pre-approve job generated. As in this snippet from the pipeline documentation page:  
```yaml
      - terraform/plan:
          name: plan preview-us-west-2 change
          context: *context
          shell: op run --env-file op.nonprod.env -- /bin/bash -eo pipefail
          executor-image: *executor-image
          static-analysis: false
          workspace: preview-us-west-2
          tfc-local-execution-mode: true
          tfc-organization: twdps
          tfc-workspace-name: psk-aws-platform-vpc-preview-us-west-2
          before-terraform:
            - set-environment:
                instance-name: preview-us-west-2
                env-credentials: nonprod
          filters: *on-tag-main
```

The name of the pre- job created for the preview-us-west-2 instance of the preview role will be configured based on this string value. Define a custom pre- job name pattern by passing this setting a custom string. The %s in the string will be substitued with the instnace or role name.  

#### 6.1.4 `--PipePostJobName`

**default** _"apply %s change"_  

The configuration of the job name for each post-approve job generated.

#### 6.1.4 `--PipeApprovalJobName`

**default** _"approve %s changes"_  

The configuration of the job name for each Approval job generated

#### 6.1.4 `--EnvFilesCreate`

**default** _true_  

Should environment files be generated automatically with creating a pipeline.  

#### 6.1.4 `--EnvFilesPath`

**default** _environments_  

"Path/to" folder containing env files.  

#### 6.1.4 `--DefaultsFileName`

**default** _default_  

Name of the env file containg default settings. This is a required file.  

#### 6.1.4 `--EnvFilesExt`

**default** _.json_  

Extension to the default, pipeline, role, or instance environment files maintained in the EnvFilesPath folder. Note, this is not the extention of the generated envfile. This is the extension used when searching for the env files you maintain in the environment/ folder with the default, pipielin, role, and/or instance specific env value settings.  

Suported options are .json, .yaml, .yml  

#### 6.1.4 `--EnvFilesWriteExt`

**default** _.tfvars.json_  

Extension applied to the generated role and instance envfiles used when created the pipeline.  



#### 6.2 Experimental options  

These flags are not recommended to be used.  

###### 6.2.1 `--PipeWorkflowHeading`

**default**  
```yaml
workflows:
  version: 2

  %s:
    jobs:
```
###### 6.2.1 `--PipeApprovalTemplate`

**default**  
```yaml
      - {{.jobname}}:
          type: approval
          {{.jobstobeapproved}}
          filters: {{.filter}}
```
Name of a file containing the Approval step template.  

###### 6.2.1 `--PipePath`

**default** _.circleci_  

Path to the circleci pipeline file.

###### 6.2.1 `--CircleciConfigFile`

**default** _config.yml_  

Name of thecircleci pipeline config file.



#### 6.3 Circlepipe configuration

Circlepipe supports a configuration file that can maintain custom overrides instead of using --flags. By default, circlepipe will look for this file in the pwd. For circleci pipeline runs this is the same directory where the .circleci folder is located.  

By default this file is presumed to be `.circlepipe.yaml` however you can specify a different path/file using the `--config=` flag. See the [Cobra](https://github.com/spf13/cobra) go module documentation for supported file types.  

###### 6.3.1 Generate a config from from the current configuration settings

Use `config init` command to generate a .circlepipe.yaml file with all supported values.  

###### 6.3.1 View the current configuration settings

Use the `config view` command to see all the current settings.  

Example output:
```json
{
  "circleciconfigfile": "config.yml",
  "envdefaultsfilename": "default",
  "envfilescreate": true,
  "envfilesext": ".json",
  "envfilespath": "environments",
  "envfileswriteext": ".tfvars.json",
  "pipeapprovaljobname": "approve %s changes",
  "pipeapprovaltemplate": "default",
  "pipecontrolfilename": "generate.yaml",
  "pipeisapprove": true,
  "pipeispost": true,
  "pipeispre": true,
  "pipeoutfile": "generated_config.yml",
  "pipepath": ".circleci",
  "pipepostjobname": "apply %s change",
  "pipepostroleonly": false,
  "pipeposttemplate": "post-approve.yml",
  "pipeprejobname": "plan %s change",
  "pipepreroleonly": false,
  "pipepretemplate": "pre-approve.yml",
  "pipepriorjobsrequired": true,
  "pipeworkflowheading": "default",
  "pipeworkflowname": "continuation-generated-workflow"
}
```

###### 6.3.2 Get a particular setting

Use `config get NAME` to see the current setting of a specific configuration value.  


###### 6.3.2 Set a configuration value

Use `config set NAME VALUE` to update the contents of the .circlepipe.yaml configuration file with a customization you want automatically applied to every circlepipe .  

###### 6.3.3 Use ENV variables

Each of the above --flag customizations can also be set using ENV variables. Use the following format to create an ENV variable to any flag:  

`export CIRCLEPIPE_FlagName="custom value`

Substitute the desired flag name in the format.  

#### 6.4 Other Pipeline formatting examples

###### 6.4.1 Run a job for all instances at once, in parallel

This is a common use case when running nightly (recurring) integration tests. For example, continuing to use the EKS terraform pipeline example, it is a good practice to perform a nightly integration test to watch for unexpected configuration changes. Infrastructure pipeline are uniquely suseptible to these since virtually every configuration made will have a manual, alternate means of being set.  

In this situation, I want the generated pipeline to trigger jobs targeting each instance of each role concurrently, with no approval step nor any requirement that one job run before or after another.  

To generate this type of pipeline, create an appropriate job template. Here is an example named `nightly.yml`:  
```yaml
      - integration-tests:
          name: {{.instance}} integration test
          context: *context
          instance-name: {{.instance}}
          env-credentials: {{.env_credentials}}

```
This template assumes that you have added a job to base config.yml called integration-tests that requires the defined parameters.  

Generate the pipeline by running the `create pipeline PIPELINENAME` command with the following flag settings
```
--PipeIsApprove=false \                      # no approval step is required
--PipeIsPost=false \                         # no post- jobs are required
--PipePreTemplate=nightly.yml \              # use a template called nightly.yml for the pre- jobs
--PipePreJobName="%s integration test" \     # rename the Pre- job name to a more accurate description
--PipePriorJobsRequired=false \              # no prior jobs are required by any of the pre- jobs
--PipeWorkflowName="nightly-scheduled-instance-integration-test"  # use a more description workflow name
```

For the example control file used in the FULL example are of the pipeline documentation page, this is how the workflow section of the generated_config.yml would look:
```yaml

workflows:
  version: 2

  nightly-scheduled-instance-integration-test:
    jobs:
      - integration-tests:
          name: preview-us-west-2 integration test
          context: *context
          instance-name: preview-us-west-2
          env-credentials: nonprod

      - integration-tests:
          name: preview-us-east-2 integration test
          context: *context
          instance-name: preview-us-east-2
          env-credentials: nonprod

      - integration-tests:
          name: preview-eu-west-1 integration test
          context: *context
          instance-name: preview-eu-west-1
          env-credentials: nonprod

      - integration-tests:
          name: preview-eu-central-1 integration test
          context: *context
          instance-name: preview-eu-central-1
          env-credentials: nonprod

      - integration-tests:
          name: nonprod-us-west-2 integration test
          context: *context
          instance-name: nonprod-us-west-2
          env-credentials: nonprod

      - integration-tests:
          name: nonprod-us-east-2 integration test
          context: *context
          instance-name: nonprod-us-east-2
          env-credentials: nonprod

      - integration-tests:
          name: nonprod-eu-west-1 integration test
          context: *context
          instance-name: nonprod-eu-west-1
          env-credentials: nonprod

      - integration-tests:
          name: nonprod-eu-central-1 integration test
          context: *context
          instance-name: nonprod-eu-central-1
          env-credentials: nonprod

      - integration-tests:
          name: prod-us-west-2 integration test
          context: *context
          instance-name: prod-us-west-2
          env-credentials: prod

      - integration-tests:
          name: prod-us-east-2 integration test
          context: *context
          instance-name: prod-us-east-2
          env-credentials: prod

      - integration-tests:
          name: prod-eu-west-1 integration test
          context: *context
          instance-name: prod-eu-west-1
          env-credentials: prod

      - integration-tests:
          name: prod-eu-central-1 integration test
          context: *context
          instance-name: prod-eu-central-1
          env-credentials: prod

```

###### 6.4.1 Run jobs at the role-level only  


<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
