<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## 4 envfiles

When deploying infrastructure or applications, there is usually a set of configuration values that go along with each discrete deployment. If you have three environments in your deployment pipeline then you typically have three *.env files of some form to go along with each deployment.  

While some settings will be unique to the environment, others may be common across all. That sort of duplication is usually fairly manageable by combining shared and unique config files at the time of deployment. But not all the frameworks and/or tools that are used together in a deployment necessarily have the same level of native support for this issue.  

Circlepipe provides a mechnism for automating this process as part of the pipeline creation. It is currently limited to json/yaml formatted data. This provides a lot of flexability and use of such env files is fairly common in infrastructure pipelines in particular.

#### 4.1 Environments folder

Circlepipe expects all of the envfiles to be in a dedicated folder. Envfiles it creates will also be written to this folder. By default it will refer to `environments/`, though this is customizable.  

#### 4.2 Environment files

Which files will be actively maintained in this environments folder (rather than generated) is up to you, though there is one required file. The names of the files that do exist must match one or more of the following:
- default (_required_)
- pipeline
- role
- instance

For example, If I have a pipeline defined in the pipeline control file as follows:
```yaml
---
sandbox:
  filter: "*on-push-main"
  deploy:
    - sbxdev
    - sbxqa
    - sbxmapi
  roles:
    sbxdev:
      deploy:
        - sbxdev-us-west-2
        - sbxdev-eu-west-1
      instances:
        sbxdev-us-west-2:
          from_generate_aws_region: us-west-2
          from_generate_aws_account_id: 'sbx10100000000'
        sbxdev-eu-west-1:
          from_generate_aws_region: eu-west-1
          from_generate_aws_account_id: 'sbx10100000000'
    sbxqa:
      deploy:
        - sbxqa-us-west-2
        - sbxqa-eu-west-1
      instances:
        sbxqa-us-west-2:
          from_generate_aws_region: us-west-2
          from_generate_aws_account_id: 'sqa10100000000'
        sbxqa-eu-west-1:
          from_generate_aws_region: eu-west-1
          from_generate_aws_account_id: 'sqa10100000000'
    sbxmapi:
      deploy:
        - sbxmapi-us-west-2
      instances:
        sbxmapi-us-west-2:
          from_generate_aws_region: us-west-2
          from_generate_aws_account_id: 'mapi10100000000'
```

Then I can choose to include one or more of the following environments files in the environment folder:
- sandbox.json
- sbxdev.json
- sbxqa.json
- sbxmapi.json
- sbxdev-us-west-2.json
- sbxdev-eu-west-1.json
- sbxqa-us-west-2.json
- sbxqa-eu-west-1.json
- sbxmapi-us-west-2.json

What these represent are environmental configuration values that are defined at the assoicated level.  

In any deployment pipeline there may be settings that will be the same across all deployments; some that are unique to a role though the same for all instance; and some that are unique to the instance.  

Ex: Continuing with concepts from our earlier examples, here is a snippet from a typical terraform EKS deployment:
```json
{
  "aws_assume_role": "PSKPlatformEksBaseRole",
  "aws_region": "us-east-2",
  "aws_account_id": "101000000000",
  "cluster_enabled_log_types": [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler"
  ],
  "cluster_log_retention": "30",
  "cluster_version": "1.27",
  "default_node_group_ami_type": "AL2_x86_64",
  "default_node_group_capacity_type": "SPOT",
  "default_node_group_disk_size": "50",
  "default_node_group_instance_types": [
    "t2.2xlarge",
    "t3.2xlarge",
    "t3a.2xlarge",
    "m5n.2xlarge",
    "m5.2xlarge",
    "m4.2xlarge"
  ],
  "default_node_group_max_size": "5",
  "default_node_group_min_size": "3",
  "default_node_group_desired_size": "3",
  "default_node_group_name": "group-a",
  "default_node_group_platform": "linux",
  ...
}
```
The example pipeline from above represents a typical development release pipeline for an engineering platform team's non-customer facing environments. In other words, these are that team's "dev" environments. There would also be something like a "release" pipeline in the control file that would generate the deployment pipeline for successful changes to be pushed to all customer-facing environments.  

In this situation, you could imagine that from the above env snippet the following settings would be the same across all pipelines, roles, and instances.  Therefore, these setting would be defined in `default.json`. (This default filename is customizable)
```json
{
  "aws_assume_role": "PSKPlatformEksBaseRole",
  "cluster_version": "1.27",
  "cluster_enabled_log_types": [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler"
  ],
  "default_node_group_ami_type": "AL2_x86_64",
  "default_node_group_disk_size": "50",
  "default_node_group_platform": "linux",
  ...
}
```
The following settings would likely be the same across all the "dev" roles and instances, which means the entire "sandbox" pipline. Therefore these settings would go in `sandbox.json`.
```json
{
  "cluster_log_retention": "7",
  "default_node_group_capacity_type": "SPOT",
  "default_node_group_instance_types": [
    "t2.2xlarge",
    "t3.2xlarge",
    "t3a.2xlarge",
    "m5n.2xlarge",
    "m5.2xlarge",
    "m4.2xlarge"
  ],
  "default_node_group_max_size": "5",
  "default_node_group_min_size": "3",
  "default_node_group_desired_size": "3",
  "default_node_group_name": "group-a",
  ...
}
```
In this example, 7 days of log retention is ample for the sandbox-dev and sandbox-qa roles, but let's assume that the sandbox management cluster needs 14 days of retention. This setting would go in `sbxmapi.json`.
```json
{
  "cluster_log_retention": "14",
  ...
}
```
And finally, there can be settings that are specific to a role instance. In this example the specific aws region is unique to the cluster. And we could, theoretically, put such a config in each of an associated instance specific envfile (`sbxdev-us-west-2.json`, etc). However, I will describe in the section below why this may not be necessarily the case.  

#### 4.3 `create envfiles pipelineName`


The `create envfiles` command ingests these files and merges them together logically to create the resulting per-instance file to be used by the pipeline for the given instance.  By default circlepipe will use the extention `.tfvars.json` with the newly created file, however you can specify the resulting file extension with json and yaml being supported. Generally it is a good practice to use something that doesn't reesult in exactly the same fqfn as the persistent envfiles.

The pipeline control file describes the order and conditional inclusion. In other words, circlepipe will iterate over the definition in the control file of a specified pipeline. It merges the files together with the following logic:

```
with pipeline

    for each role in roles
        settings = default.json
        settings += pipeline.json
        settings += role.json

        for each instance in instances
            settings += instance.json
            settings += all key:value pairs from the control file for the particular instance
            write instance.tfvars.json (See customization.md for default and how to override)
        end

    end
```
<sup>(psuedocode)</sup>

Note, since all key-value pairs that are part of an instance definition in the pipeline control file are also merged into the instance envfile I do not need to duplicate that information within any repos using the control file. So in the case of the aws_region and aws_account_id, those are picked up from the control file and we do not need to specify them in our managed instance settings.  

#### 4.4 Integration with `create pipeline` command

While chiefly the resulting environment files are a dependency for actions that take place during the actual deploy pipeline jobs, they also come into play when defining the circleci pipeline templates used by the `create pipeline` command.  

The `create pipeline` command is used to generate the generate_config.yml called by the continutation orb. This command will first, by default, run the `create envfiles` command to generate the instance envfiles. You can change this default action if you want to fully manage the instance envfiles yourself.  

Any value that is in these resulting instance envfiles may be referenced from within a circleci pipeline template file.  

See [pipelines](./pipelines.md) documentation for details usage.  

#### 4.45 Additional requirements

**caching the instance envfiles**  
As noted above, while you can reference the contents of the instance envfiles within your circleci pipeline templates, these primarily exist to be consumed by the deploy jobs themselves in the resulting pipeline. Therefore, you must store the files in a cirecleci cache layer for them to be available for jobs within the generated pipelines.  

See [pipelines](./pipelines.md) for a detailed description of how to add this to your pipeline.

**install necessary tools**  
The most convient way to utilize the circleci continuation capability is with their continuation orb wherein they provide an executor with the appropriate code installed. However, keep in mind that if the continuation executor is used then an additional requirement arises.  

The circlepipe tool will not be present on the continuation executor. Nor will the other likely required tools such a tool used to fetch the pipeline control file. In the job that launches the continuation pipeline you will need to install circlepipe and other such requirements.  

Both json and yaml file formats are supported for the enduring env files and the resulting generated files. See [customization](./customization.md) for detailed instructions.  


<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
