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

Suppose you are part of an Engineering Platform product team where the compute featuer is built around kubernetes. For a quality, enterprise level platform this results in a k8s release pipeline that would look something like this:

<div align="center">
		<img alt="basic pipeline" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/circlepipe/main/basic-ptp.png?sanitize=true" />
    <br />
</div>

The platform is designed to support a global product offering. In this case this means that, as a developer-user of the platform, your production namespace (along with your other environments) actually represents namespaces in k8s clusters located in mulitiple regions around the world. As the platform product owner, the typical path-to-production for platform features

<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
