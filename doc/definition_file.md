<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## 3 multi-environment pipeline defintion file

The pipeline definition file defines all the pipelines you want to be able to generate from within the base pipeline of a repository. The definition of a pipeline is built around an ordered list of the environments you want to deploy and the list of all the instances of the environment.  

Suppose you are part of an Engineering Platform product team where the compute feature is built around kubernetes. For a quality, enterprise level platform this results in a k8s release pipeline that would look something like this:

![basic path to production](images/basic-ptp.png)

Developer-users (customers) of the platform will have namespaces in the preview, nonprod, and prod clusters. However, when evolving to support multi-regional compute definition, the platform team will now need to have preview, nonprod, and prod clusters in multiple regions around the world. Hence a kuberneset change to _preview_ now means a change to many clusters rather than a single instace.

`git push`  is still expected to result in deployments to all of the product test environments and `git tag` still deploys changes to all the customer-facing clusters. Depending on the trigger, there are two different pipelines that could need to be generated.  

Using yaml, let's start be defining the two possible pipelines and the associated triggers[^1]:
```yaml
---
test:
  filter: "*on-push-main"
release:
  filter: "*on-tag-main"
```

These anchors would need to be in your base config.yml. Naturally you can define these trigers however you want in the context of your pipeline.  

The platform is designed to support a global product offering. In this case this means that, as a developer-user of the platform, your production namespace (along with your other environments) actually represents namespaces in k8s clusters located in mulitiple regions around the world. As the platform product owner, the typical path-to-production for platform features



**footnotes**

[^1]:In the example below, the "\*on-push-main" and "\*on-tag-main" notation refers to a specific anchor definition in yaml for circleci that looks like this
```yaml
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
<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
