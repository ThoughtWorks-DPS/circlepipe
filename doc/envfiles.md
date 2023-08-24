<div align="center">
	<p>
		<img alt="Thoughtworks Logo" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/thoughtworks_flamingo_wave.png?sanitize=true" width=200 />
    <br />
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h4>circlepipe documentation</h4>
</div>
<br />

## envfiles

    - save_cache:
        name: |
          circlepipe can generate ENV setting files automatically as an optional feature. If
          you use these feature you must persist the generated files between workflows.
          (Read detailed instructions on how to use envfiles in the circlepipe documentation
          under the toc entry 'Create Envfiles'.)
        key: circlepipe-{{ .Revision }}
        paths:
          - environments/

<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
