package cmd

import (
	"testing"
	"errors"
)

func TestParameterValidation(t *testing.T) {
	tests := []struct {
		job string
		name string
		expected error
	}{
		{ "pipeworkflowname", "a valid workflow or job name", nil },
		{ "pipeworkflowname", "an !nvalid workflow or job name", errors.New("circlepipe set error: workflow name too long or contains invalid characters") },
		{ "PipePreJobName", "a valid workflow or job name", nil },
		{ "PipePostJobName", "a valid workflow or job name", nil },
		{ "PipeApprovalJobName", "a valid workflow or job name", nil },
		{ "pipeisapprove", "true", nil},
		{ "pipeisapprove", "false", nil},
		{ "pipeisapprove", "invalid", errors.New("circlepipe set error: boolean key value setting must be true or false") },
		{ "pipepriorjobsrequired", "true", nil},
		{ "pipeIsPre", "true", nil},
		{ "PipeIsPost", "true", nil},
		{ "PipePreRoleOnly", "true", nil},
		{ "PipePostRoleOnly", "true", nil},
		{ "envfilescreate", "true", nil},
		{ "pipecreateapprovalstep", "true", nil},
		{ "piperoleonly", "true", nil},
		{ "PipePreTemplate", "set_test.go", nil},
		{ "PipePreTemplate", "I:!valid", errors.New("stat I:!valid: no such file or directory")},
		{ "PipePreTemplate", "nonexistant.file", errors.New("stat nonexistant.file: no such file or directory")},
		{ "PipePostTemplate", "set_test.go", nil},
		{ "envfilespath", "set_test.go", nil},
		{ "envdefaultsfilename", "set_test.go", nil},
		{ "pipepath", "set_test.go", nil},
		{ "pipecontrolfilename", "set_test.go", nil},
		{ "PipeApprovalTemplate", "set_test.go", nil},
		{ "pipeoutfile", "filename.txt", nil },
		{ "pipeoutfile", "file:name.txt", errors.New("circlepipe set error: invalid outfile name") },
		{ "EnvFilesExt", ".txt", nil },
		{ "EnvFilesExt", "txt", errors.New("circlepipe set error: extensions must start with '.'") },
		{ "EnvFilesWriteExt", ".txt", nil },
	}


	for _, test := range tests {
		result := validateSetting([]string{test.job, test.name})
		if result == nil && test.expected != nil {
			t.Errorf("For %s = %s, expected error %v, but got nil", test.job, test.name, test.expected)
		} else if result != nil && (test.expected == nil || result.Error() != test.expected.Error()) {
			t.Errorf("For %s = %s, expected error '%v', but got '%v'", test.job, test.name, test.expected, result)
		}
	}

}
