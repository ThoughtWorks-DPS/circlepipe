<div align="center">
	<p>
		<img alt="DPS Title" src="https://raw.githubusercontent.com/ThoughtWorks-DPS/static/master/EMPCPlatformStarterKitsImage.png?sanitize=true" width=350/>
	</p>
  <h2>circlepipe</h2>
  <br />
  <h3>documentation</h3>
</div>
<br />

## installation

Releases are found here.

<div class="tabbed-callout">
  <input type="radio" name="tab" id="tab1" checked>
  <input type="radio" name="tab" id="tab2">
  <input type="radio" name="tab" id="tab3">

  <div class="tabs">
    <label for="tab1">Tab 1</label>
    <label for="tab2">Tab 2</label>
    <label for="tab3">Tab 3</label>
  </div>

  <div class="content">
    <div class="tab-content">
      Content for Tab 1
    </div>
    <div class="tab-content">
      Content for Tab 2
    </div>
    <div class="tab-content">
      Content for Tab 3
    </div>
  </div>
</div>

<style>
.tabbed-callout .tabs {
  display: flex;
}

.tabbed-callout .tabs label {
  flex: 1;
  text-align: center;
  padding: 10px;
  cursor: pointer;
}

.tabbed-callout .content {
  border: 1px solid #ccc;
  padding: 10px;
}

.tabbed-callout .tab-content {
  display: none;
}

.tabbed-callout input[type="radio"]:checked + .content .tab-content:nth-child(1),
.tabbed-callout input[type="radio"]:checked + .content .tab-content:nth-child(2),
.tabbed-callout input[type="radio"]:checked + .content .tab-content:nth-child(3) {
  display: block;
}
</style>


curl -SLO https://github.com/ThoughtWorks-DPS/circlepipe/releases/latest/download/circlepipe_Linux_amd64.tar.gz && \
tar -xzf opw_Linux_x86_64.tar.gz && \
sudo mv opw /usr/local/bin/opw && \


<hr>  

[<kbd> <br> Back <br> </kbd>](./table_of_contents.md)
