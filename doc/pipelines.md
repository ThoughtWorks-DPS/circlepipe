<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## 5 pipelines

There are four steps involved in successfully generating a pipeline using circlepipe.  

- Create a pipeline control definition file that contains definitions for all the pipelines you want to be able to generate from the initially triggered pipeline.
- Modify the initially triggered pipeline with the necessary launch jobs and workflow instructions that will use circlepipe to generate a new pipeline and then launch it using the CircleCI continuation orb.
- Create any needed [envfiles](./envfiles.md)
- Create the circleci pipeline job templates used by circlepipe in generating the continuation pipeline.

#### 5.1 Create a pipeline control definition file

See [multi-environment pipeline defintion file](./definition_file.md).  

#### 5.2 Add configuration to the initially triggered pipeline

##### 5.2.1 Add the continuation orb and `setup:` directive to config.yml

Typically, this goes near the top:  
```yaml
version: 2.1

setup: true

orbs:
  continuation: circleci/continuation@0.4.0
  ...
```

##### 5.2.2 Add any jobs, commands, or other pipeline configuration used by the generated pipelines

These may or may not also be commands used within the initially triggered pipeline. An example of a shared command would be a general environment setup command. But it is also common to need jobs or commands to be used only within the generated pipeline.  

The pipeline generated by the `create pipeline` command is a copy of the initial config.yml with the setup:true directive and all existing workflows removed and populated with a new, generated workflow.  

In addition to jobs and command this means that any anchors or pipeline parameters defined in the config.yaml can be used by the generated pipeline.  Filters are an example of an anchor that is required by a pipeline definition in the control file in order to properly construct the pipeline.  

```yaml

parameters:
  continuation-params:
    description: parameters to be passed to continuation orb
    type: string
    default: ""

globals:
  - &context empc-lab

on-push-main: &on-push-main
  branches:
    only: /main/
  tags:
    ignore: /.*/

on-tag-main: &on-tag-main
  branches:
    ignore: /.*/
  tags:
    only: /.*/
```

##### 5.2.3 Add a job that uses the continuation orb to launch a new pipeline

In the jobs: section of the config.yml file, create a job for launching a generated pipeline using the continuation orb.  

```yaml
jobs:

  ...

  launch-dynamic-pipeline:
    parameters:
      pipeline-name:
        description: name of the pipeline definition from the pipeline control file
        type: string
    executor: continuation/default
    steps:
      - checkout
      - setup_remote_docker
      - run:
          name: |
            When using the continuation orb executor for launching a generated pipeline, naturally circelpipe
            nor any other custom tool you need for the "create pipeline" step will be available on the
            executor.

            Before generating the pipeline then, install these tools. In this example both circlepipe and
            a secrets manage tool are installed since the example assumes the official copy of the pipeline
            control file is maintained within the secrets management service.  
          command: |
            curl -SLO https://github.com/ThoughtWorks-DPS/circlepipe/releases/latest/download/circlepipe_Linux_amd64.tar.gz
            tar -xzf circlepipe_Linux_amd64.tar.gz
            sudo mv circlepipe /usr/local/bin/circlepipe
            curl -L https://cache.agilebits.com/dist/1P/op2/pkg/v2.18.0-beta.01/op_linux_amd64_v2.18.0-beta.01.zip -o op.zip
            unzip -o op.zip && sudo mv op /usr/local/bin/op
      - save_cache:
          name: |
            You must persist the generated envfiles files between workflows. (Read detailed instructions
            on how to use envfiles in the circlepipe documentation under the toc entry 'Create Envfiles'.)
          key: circlepipe-{{ .Revision }}-<< parameters.pipeline-name >>
          paths:
            - environments/
      - run:
          name: Generate the new pipeline
          command: circlepipe create pipeline << parameters.pipeline-name >>
      - continuation/continue:
          name: |
            The generated pipeline has been created and you can now use the continuation orb to launch it.
            Keep in mind that a continuation step can only be called from the initial launched pipeline.
            You cannot include setup:true in the generated pipeline in generate additional pipelines.
          parameters: |
            { "continuation-params": "<< pipeline.parameters.continuation-params >>", "pipeline-name": "<<parameters.pipeline-name >>"" }
          configuration_path: .circleci/generated_config.yml
```

##### 5.2.4 Include the "launch continuation job" in the pipeline workflow

In this example, the workflow triggered by git-push first performs any CI type jobs prior to generating and launching the "sandbox" deployment pipeline. Those jobs, assuming they are successful, only run with a new code push. To trigger a release to production, the commit is tagged, this results in the "release" pipeline being generated and launched.  

```yaml
workflows:
  version: 2

  development build:
    jobs:

      ...

      - static-analysis:
          # You can include jobs that you want to run before continuing with a new pipeline.
          # A common pattern is to perform all the static code analysis, unit tests and coverage
          # reporting and other similar "only at build time" tests, and then launch a deployment
          # to move the code on to production
          filters: *on-push-main

      - launch-dynamic-pipeline:
          name: generate sandbox pipeline
          context: *context
          pipeline-name: sandbox
          requires:
            - static-analysis
          filters: *on-push-main

  release candidate:
    jobs:
      - launch-dynamic-pipeline:
          name: generate release pipeline
          context: *context
          pipeline-name: release
          filters: *on-tag-main
```



#### 5.3 Create pipeline, role, and instance envfiles as needed

See [creating envfiles](./envfiles.md).

#### 5.4 Create the circleci pipeline job templates

Circlepipe generates a deployment pipeline based on _templates_ you provide. These templates are portions of a circleci workflow than include notations where information from the instance environment files and the pipeline control file can be injected.  

You can define two types of templates.  

1. A "pre" job to be called for each instance of a role, or optionally just the role, prior to an optional Approve step.
2. A "post" job to be called for each instance of a role, or optionally just the role, after an Approval step.

**pre-approval.yml**  

The default name for the template file containing jobs that come before an optional approval is `pre-approve.yml` though this can be customized. By default, cirlepipe will look for the template `.circleci/` folder, which also is customizable.  

The template can contain any valid circleci workflow items, however, in most instances it will contain a single job reference. Continuing with the EKS terraform deployment example, the pre-appove template could look like this: (note, the spacing is required)
```yaml
      - terraform/plan:
          name: plan {{.instance}} change
          context: *context
          shell: op run --env-file op.{{.env_credentials}}.env -- /bin/bash -eo pipefail
          executor-image: *executor-image
          static-analysis: false
          workspace: {{.instance}}
          tfc-local-execution-mode: true
          tfc-organization: twdps
          tfc-workspace-name: psk-aws-platform-vpc-{{.instance}}
          before-terraform:
            - set-environment:
                instance-name: {{.instance}}
                env-credentials: {{.env_credentials}}
          filters: {{.filter}}
          {{.priorapprovalrequired}}

```
In this example template, the job is one of jobs provided by the twdps/terraform orb. This orb contains a terraform plan job that includes a variety of optional steps and configurations. You can read more about the orb here. But in general terms, the job will perform all the necessary actions of a terraform plan action.  

You will see several template injection fields indicated by the double braces, such as {{.instance}}. A template can be configured to inject any value that appears in the current instance's env file. In addition, there are several other values added to the instance env file by circlepipe to aid in customizing the behavior of the template. These are:

- {{.filter}} = from the filter key in the pipeline control file for the defined pipeline.
- {{.role}} = The role that includes the current instance.
- {{.instance}} = The instance name.
- {{.envfilespath}} = the folder which contains all the envfiles.
- {{.priorapprovalrequired}} = The name of the prior approval job for which a "pre" template needs to include a "requires:"
- {{.jobname}} = the name to use for the approval job when injecting the approval step into the pipeline config.
- {{.jobstobeapproved}} = The list of jobs which must all be completed prior to an approval step.
- {{.approvaljobname}} = The name of the prior approval job for which a "post" template needs to include a "requires:"
- {{.lastjob}} = the job immediately before the job being templated

**approve template**  
The approval step template looks like this by default:
```yaml
      - {{.jobname}}:
          type: approval
          {{.jobstobeapproved}}
          filters: {{.filter}}

```
While it is possible to customize this layout by providing an alternative yaml file, this would be unusual and is basically an experimental option.  

**post-approve.yml**  
Continuing with the above EKS terraform example, this is what the post-approval job would look like:
```yaml
      - terraform/apply:
          name: apply {{.instance}} change
          context: *context
          shell: op run --env-file op.{{.env_credentials}}.env -- /bin/bash -eo pipefail
          workspace: {{.instance}}
          before-terraform:
            - set-environment:
                instance-name: {{.instance}}
                env-credentials: {{.env_credentials}}
          after-terraform:
            - run-inspec-tests:
                instance-name: {{.instance}}
          requires:
            - {{.approvaljobname}}
          filters: {{.filter}}

```
This also is making use of a job defined by the twdps/terraform orb. Notice that there are `before-` and `after-` workflow hooks that allow any custom command to be added to the job. In this case, before running terraform apply the needed pipeline environment settings are configured. Then after the apply step commands are called that perform integration tests to validate the successful infrastructure changes.  

Look at this set of files that make up a full set of files in a complete example of the sample EKS infrastructure pipeline referred to repeatedly above.  

| file | description |
|------|-------------|
| [config.yml](../test/test_pipeline_terraform/result/config.yml) | the initial triggered pipeline config.yml |
| [pre-approve.yml](../test/test_pipeline_terraform/result/pre-approve.yml) | pre-approve template for terraform style workflow |
| [post-approve.yml](../test/test_pipeline_terraform/result/post-approve.yml) | post-approve template for terraform style workflow |
| [generate.yaml](../test/test_pipeline_terraform/result/generate.yaml) | pipeline control file |
| [default.json](../test/test_pipeline_terraform/result/default.json) | default values for all pipelines, roles, instances |
| [release.json](../test/test_pipeline_terraform/result/release.json) | overrides or additions  for the release pipeline |
| [nonprod.json](../test/test_pipeline_terraform/result/nonprod.json) | overrides or additions for the nonprod role |
| [prod.json](../test/test_pipeline_terraform/result/prod.json) | overrides or additions for the prod role |
| [generated_config.yml](../test/test_pipeline_terraform/result/generated_config.yml) | the resulting generated pipeline to be launched by the continuation orb |

Here is a screen shot of the circleci visualization of the pipeline.  

<div align="center">
	<p>
		<img alt="example pipeline" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/circlepipe/main/doc/images/example_pipeline.png" />
	</p>
</div>  

#### 5.5 Other pipeline workflow options

Circlepipe has several customization options the enable a number of other workflow layout beyond just "pre + approval + post" by-instance demonstrated in thie above example. See the [customization](./customization.md) page for additional details and examples.  

<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
