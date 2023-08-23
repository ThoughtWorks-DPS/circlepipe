<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <br />
  <h3>circlepipe</h3>
    <a href="https://app.circleci.com/pipelines/github/ThoughtWorks-DPS/circlepipe"><img src="https://dl.circleci.com/status-badge/img/gh/ThoughtWorks-DPS/circlepipe/tree/main.svg?style=shield"></a> <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-blue.svg"></a>
</div>
<br />

Opinionated generation of CircleCI [dynamic configuration](https://circleci.com/docs/using-dynamic-configuration/) _continuation_ pipelines[^1].

## Purpose

The primary motivation for the development of this tool is in generating pipelines for deploying multi-regional architectures. Once you get to that scale, many infrastructure or application pipelines need to be capable of dynamically responding to added or substracted regions and environments.  

Imagine you support a globally-available microservice architecure that runs on EKS. This means you must create multiple EKS instances, spanning the supported regions. What if _Production_ constitutes eks clusters in 3 localities around the world, with additionally 2 regions per locality where traffic is geo-location routed to the closest in proximity? This can result in six (6) clusters, that together actually make up just one environment - Production. Using CircleCI, and assuming you have cleanly abstracted the actual deploy/test/etc steps into maintained orbs, a low-complexity solution is to dynamically generate the 'deployment' pipeline at runtime.

For example, suppose I have a terraform orb that incorporates all of my requirements for performing `terraform plan` and `apply` like this:
```
      - terraform/plan:
          name: plan nonprod-us-west-2 change
          context: *context
          shell: op run --env-file op.nonprod.env -- /bin/bash -eo pipefail
          static-analysis: false
          workspace: nonprod-us-west-2
          before-terraform:
            - set-environment:
                instance-name: nonprod-us-west-2
          filters: *on-tag-main

      - approve nonprod changes:
          type: approval
          requires:
            - plan nonprod-us-west-2 change
          filters: *on-tag-main

      - terraform/apply:
          name: apply nonprod-us-west-2 change
          context: *context
          shell: op run --env-file op.nonprod.env -- /bin/bash -eo pipefail
          workspace: nonprod-us-west-2
          before-terraform:
            - set-environment:
                instance-name: nonprod-us-west-2
          after-terraform:
            - run-tests:
                instance-name: nonprod-us-west-2
          requires:
            - approve nonprod changes
          filters: *on-tag-main
```

When I move to a multi-region architecture, wouldn't it be convenient if I could simply create a definition of my EKS environments sort of like this:
```json
{
  "nonprod": [
    "nonprod-us-west-2",
    "nonprod-us-east-2",
    "nonprod-eu-west-1",
    "nonprod-eu-central-1"
  ]
}
```
And then by passing this definition to a _generate circleci pipeline_ tool I would get a deployment pipeline like I use above but with plan and apply jobs for all the environments in my definition. In a nutshell, that is what `circlepipe` is designed to do.

Regardless of how many steps are needed to create a deployment, what is manageable when each environment amounts to a single deployment can quickly become unmanageable where a given environment is really many locations (and repeatedly increases).  

Read the full usage documentation [here](doc/table_of_contents.md).  

## installation


## Contributing

We encourage [issues](https://github.com/ThoughtWorks-DPS/circlepipe/issues) and [pull requests](https://github.com/ThoughtWorks-DPS/circlepipe/pulls) against this repository. In order to value your time, here are some things to keep in mind:  

1. This tool is in active development and both small and radical changes might be in the works at any given time. We will try to keep things like that visible in the Issues so always take a look there before deciding on change.  
2. The current release is still 0.*. Which is not to say the tool isn't in active use. But rather that there has not been sufficient time for truly exhaustive testing. We document and address bugs when we find them but don't be surprised if you find some as well. :smile:

[^1]:_note. This project originally started with [circlecigen](), which is now deprecated and replaced with this tool._  
