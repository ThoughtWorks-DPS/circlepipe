package cmd

import (

	"os"
	"fmt"
	//"io"

	"testing"
	//"errors"
	"encoding/json"
	"path/filepath"
	"gopkg.in/yaml.v2"
	"sort"
	//"github.com/spf13/viper"
	"strings"
)

func anyStaticEnvFileValues(fn string) map[string]interface{} {
	// if the file does exist, do not fail but return an empty map
	if _, err := os.Stat(envStaticFile(fn)); err != nil { return make(map[string]interface{}) }
	filecontents, err :=envFileValues(envStaticFile(fn))
	exitOnError(err)
	return filecontents
}

func anyEnvFileValues(fn string) map[string]interface{} {
	// if the file does exist, do not fail but return an empty map
	if _, err := os.Stat(envFile(fn)); err != nil { return make(map[string]interface{}) }
	filecontents, err :=envFileValues(envFile(fn))
	exitOnError(err)
	return filecontents
}

func envFileValues(fqfn string) (map[string]interface{}, error) {
	var data = make(map[string]interface{})
	var err error
	if fileContents, readerr := os.ReadFile(fqfn); readerr == nil {
		switch filepath.Ext(fqfn) {
		case ".json":
			err = json.Unmarshal(fileContents, &data)
		case ".yaml", ".yml":
			err = yaml.Unmarshal(fileContents, &data)
		default:
			err = fmt.Errorf(" %s is not a supported file type", filepath.Ext(fqfn))
		}
		return data, err
	} else {
		return data, readerr
	}
}

func textFromFile(fromfile string) string {
	content, err := os.ReadFile(fromfile)
	exitOnError(err)

	return string(content)
}

func writeEnvFileValues(datamap map[string]interface{}, outfile string) {
	var data []byte
	var err error
	fqfn := envFile(outfile)

	switch filepath.Ext(fqfn) {
	case ".json":
		data, err = json.MarshalIndent(datamap, "", "  ")
		exitOnError(err)
	case ".yaml", ".yml":
		data, err = yaml.Marshal(datamap)
		exitOnError(err)
	default:
		exitOnError(fmt.Errorf(" %s is not a supported file type", filepath.Ext(fqfn)))
	}
	exitOnError(os.WriteFile(fqfn, data, 0644))
}

func getFilenames(dir string) ([]string, error) {
	files := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func compareFile(t *testing.T, fixture, result string) (bool, error) {
	resultContent, err := os.ReadFile(result)
	if err != nil {
		t.Logf("Error reading source file %s: %v", result, err)
		return false, err
	}

	fixtureContent, err := os.ReadFile(fixture)
	if err != nil {
		t.Logf("Error reading source file %s: %v", fixture, err)
		return false, err
	}

	if string(resultContent) != string(fixtureContent) {
		t.Errorf("Result does not equal fixture for %s", result)
		return false, err
	}
	return true, nil
}

func compareJsonFiles(fixture, result string) (bool, error) {
	var resultData, fixtureData map[string]interface{}

	// read files
	resultfile, err := fileContents(result)
	exitOnError(err)

	fixturefile, err := fileContents(fixture)
	exitOnError(err)

	// marshall to json
	exitOnError(json.Unmarshal(resultfile, &resultData))
	exitOnError(json.Unmarshal(fixturefile, &fixtureData))

	// Convert the maps back to canonical JSON representations
	canonicalResult, err := marshalCanonicalJSON(resultData)
	exitOnError(err)

	canonicalFixture, err := marshalCanonicalJSON(fixtureData)
	exitOnError(err)

	// Compare the canonical JSON representations
	if canonicalResult == canonicalFixture {
		return true, nil
	}
	return false, fmt.Errorf("%s does not match", result)
}

func fileContents(fqfn string) ([]byte, error) {
	result, err := os.ReadFile(fqfn)
	if err != nil {
		return nil, fmt.Errorf("error reading %s/%v", fqfn, err)
	}
	return result, nil
}

func marshalCanonicalJSON(data map[string]interface{}) (string, error) {
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sortedData []string
	for _, key := range keys {
		value, err := json.Marshal(data[key])
		if err != nil {
			return "", err
		}
		sortedData = append(sortedData, fmt.Sprintf(`"%s":%s`, key, value))
	}

	return "{" + strings.Join(sortedData, ",") + "}", nil
}

func compareYamlFiles(fixture, result string) (bool, error) {
	var resultData, fixtureData map[string]interface{}

	// read files
	resultfile, err := fileContents(result)
	exitOnError(err)

	fixturefile, err := fileContents(fixture)
	exitOnError(err)

	// marshall to json
	exitOnError(yaml.Unmarshal(resultfile, &resultData))
	exitOnError(yaml.Unmarshal(fixturefile, &fixtureData))

	// Convert the maps back to canonical YAML representations
	canonicalResult, err := marshalCanonicalYAML(resultData)
	exitOnError(err)

	canonicalFixture, err := marshalCanonicalYAML(fixtureData)
	exitOnError(err)

	// Compare the canonical YAML representations
	if canonicalResult == canonicalFixture {
		return true, nil
	}
	return false, fmt.Errorf("%s does not match", result)
}

func marshalCanonicalYAML(data map[string]interface{}) (string, error) {
	var keys []string

	for key := range data {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var sortedData []string
	for _, key := range keys {
		value, err := yaml.Marshal(data[key])
		if err != nil {
			return "", err
		}
		sortedData = append(sortedData, fmt.Sprintf(`"%s": %s`, key, value))
	}

	return strings.Join(sortedData, "\n"), nil
}
