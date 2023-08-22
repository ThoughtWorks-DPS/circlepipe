package cmd

import (
	"os"
	"os/exec"
	"testing"
)

func TestInitCmd(t *testing.T) {
	// write default config file to test folder
	myCmd := exec.Command("go", "run", "main.go", "config", "init", "--config=test/test_config_init/result/.circlepipe.yaml")
	myCmd.Dir = "../"
	_, err := myCmd.CombinedOutput()
	//t.Log(string(output))
	if err != nil {
		t.Log(err)
	}

	// compare to known good result
	result, err := compareFile(t, "../test/test_config_init/fixture/.circlepipe.yaml", "../test/test_config_init/result/.circlepipe.yaml")
	os.Remove("../test/test_config_init/result/.circlepipe.yaml")
	if err != nil {
		t.Log(result)
		t.Error(err)
	}
}
