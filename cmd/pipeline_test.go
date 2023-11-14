package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPipelineCmd(t *testing.T) {
	tests := []struct {
		testtype string
	}{
		{"terraform"},
		{"all_instances"},
		{"role_only_no_approval"},
		{"role_only"},
		{"role_and_instance_post_approve"},
		{"instance_and_role_post_approve"},
	}
	var setEnvFilesPath string
	var setPipePath string
	var setConfigFileDefaultName string

	for _, test := range tests {
		setEnvFilesPath = fmt.Sprintf("--EnvFilesPath=test/test_pipeline_%s/result", test.testtype)
		setPipePath = fmt.Sprintf("--PipePath=test/test_pipeline_%s/result", test.testtype)
		setConfigFileDefaultName = fmt.Sprintf("--config=test/test_pipeline_%s/result/.circlepipe.yaml", test.testtype)
		t.Log(setEnvFilesPath)
		t.Log(setPipePath)
		t.Log(setConfigFileDefaultName)

		// generate pipeline
		myCmd := exec.Command("go", "run", "main.go", "create", "pipeline", "release", setConfigFileDefaultName, setEnvFilesPath, setPipePath)
		myCmd.Dir = "../"
		output, err := myCmd.CombinedOutput()

		if err != nil {
			t.Log(myCmd)
			t.Log(string(output))
			t.Error(err)
		}

		resultDir := fmt.Sprintf("../test/test_pipeline_%s/result", test.testtype)
		fixtureDir := fmt.Sprintf("../test/test_pipeline_%s/fixture", test.testtype)

		resultFiles, err := getFilenames(resultDir)
		if err != nil {
			t.Log(resultFiles)
			t.Fatalf("Error reading source directory: %v", err)
		}

		fixtureFiles, err := getFilenames(fixtureDir)
		if err != nil {
			t.Log(fixtureFiles)
			t.Fatalf("Error reading destination directory: %v", err)
		}

		// Compare files
		for _, filename := range resultFiles {
			var result bool
			var err error
			resultPath := filepath.Join(resultDir, filename)
			fixturePath := filepath.Join(fixtureDir, filename)

			// compare to known good result, only for generated pipeline file
			if strings.Contains(resultPath, "generated_config") {
				// validate generated pipeline
				myCmd := exec.Command("circleci", "config", "validate", resultPath)
				output, err = myCmd.CombinedOutput()
				if err != nil {
					t.Log(string(output))
					t.Error(err)
				}

				result, err = compareYamlFiles(resultPath, fixturePath)
				if err != nil {
					t.Logf("%s generated pipeline does not match fixture", test.testtype)
					t.Log(result)
					t.Error(err)
				}

			}
		}
	}
}
