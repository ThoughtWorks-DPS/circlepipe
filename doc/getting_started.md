<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## 2 getting started

Read about the CircleCI continuation orb [here](https://circleci.com/docs/using-dynamic-configuration/).  

#### 2.1 add the continuation orb and `setup:` directive to config.yml

Typically, this goes near the top:
```yaml
version: 2.1

setup: true

orbs:
  continuation: circleci/continuation@0.4.0
```

#### 2.2 define a job that uses the continuation orb to launch a new pipeline

Before calling the continue command from the continutation orb, you must first generate the new pipeline yaml file. This is the step circlepipe does for you based on the configuration you provide. In the circleci documentation the new
pipeline is called generated_config.yml and circlepipe uses that as the defautl as well, though you can name the file whatever you like. (Detailed documentation on how to setup these instructions [here](./pipelines.md).)  
```yaml
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
        name: install circlepipe on the continuation executor
        command: |
          curl -SLO https://github.com/ThoughtWorks-DPS/circlepipe/releases/latest/download/circlepipe_Linux_amd64.tar.gz
          tar -xzf circlepipe_Linux_amd64.tar.gz
          sudo mv circlepipe /usr/local/bin/circlepipe
    - save_cache:
        name: persist envfiles files between workflows
        key: circlepipe-{{ .Revision }}-<< parameters.pipeline-name >>
        paths:
          - environments/
    - run:
        name: Generate the new pipeline
        command: circlepipe create pipeline << parameters.pipeline-name >>
    - continuation/continue:
        configuration_path: .circleci/generated_config.yml
```

#### 2.3 call the continuation job from the pipeline workflow

```yaml
workflows:
  version: 2

  development build:
    jobs:
      - static-analysis:
          # You can include jobs that you want to run before continuing with a new pipeline.
          # A common pattern is to perform all the static code analysis, unit tests and coverage
          # reporting and other similar "only at build time" tests, and then launch a deployment
          # to move the code on to production
          filters: *on-push-main

      - launch-dynamic-pipeline:
          name: generate dev pipeline
          context: *context
          pipeline-name: dev
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

<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
