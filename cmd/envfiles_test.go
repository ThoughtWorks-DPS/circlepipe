package cmd

import (
	"fmt"
	"os/exec"
	"testing"
	"path/filepath"
	"strings"
)

func TestEnvfilesCmd(t *testing.T) {
	tests := []struct {
		filetype string
	}{
		{ "json"},
		{ "yaml"},
	}
	var setEnvFileExt string
	var setEnvFilesPath string
	var setEnvFilesWriteExt string
	var setPipeControlFileName string

	for _, test := range tests {
		setEnvFilesPath = fmt.Sprintf("--EnvFilesPath=test/test_envfiles_%s/result", test.filetype)
		setEnvFileExt = fmt.Sprintf("--EnvFilesExt=.%s", test.filetype)
		setEnvFilesWriteExt = fmt.Sprintf("--EnvFilesWriteExt=.tfvars.%s", test.filetype)
		setPipeControlFileName = fmt.Sprintf("--PipeControlFileName=generate.%s", test.filetype)
		t.Log(setEnvFilesPath)
		t.Log(setEnvFileExt)

		// generate envfiles
		myCmd := exec.Command("go", "run", "main.go", "create", "envfiles", "release", setEnvFilesPath, setEnvFileExt, setEnvFilesWriteExt, setPipeControlFileName)
		myCmd.Dir = "../"
		output, err := myCmd.CombinedOutput()

		if err != nil {
			t.Log(output)
			t.Error(err)
		}

		resultDir := fmt.Sprintf("../test/test_envfiles_%s/result", test.filetype)
		fixtureDir := fmt.Sprintf("../test/test_envfiles_%s/fixture", test.filetype)

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

			// compare to known good result, only for generated env files
			if strings.Contains(resultPath, "tfvars") {
				//t.Logf("compare %s to %s", resultPath, fixturePath)
				switch filepath.Ext(resultPath) {
				case ".json":
					result, err = compareJsonFiles(resultPath, fixturePath)
					if err != nil { t.Log("json error") }
				case ".yml", ".yaml":
					result, err = compareYamlFiles(resultPath, fixturePath)
					if err != nil { t.Log("yaml error") }
				default:
					t.Errorf("envfiles_test unsupported file type")
				}

				if err != nil {
					t.Log(result)
					t.Error(err)
				}
			}
		}
	}
}
