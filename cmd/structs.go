package cmd

import (
	"os"
	"fmt"
	"path/filepath"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Role struct {
	Deploy []string
	Instances map[string]interface{}
}

type Pipeline struct {
	Filter string
	Deploy  []string
	Roles  map[string]Role
}

type PipelineMap map[string]Pipeline

func (p *PipelineMap) NewFromFile(fqfn string) error {
	if fileContents, err := os.ReadFile(fqfn); err == nil {
		switch filepath.Ext(fqfn) {
		case ".json":
			if err := json.Unmarshal(fileContents, p); err != nil { return err }
		case ".yaml", ".yml":
			if err := yaml.Unmarshal(fileContents, p); err != nil { return err }
		default:
			return fmt.Errorf(" %s is not a supported file type", filepath.Ext(fqfn))
		}
	} else {
		return err
	}
	return nil
}
