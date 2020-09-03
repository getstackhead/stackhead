package jsonschema

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"

	"github.com/getstackhead/stackhead/cli/ansible"
)

// isInternalError determines if a given error message is related to the schema requirements itself
func isInternalError(errorType string) bool {
	switch errorType {
	case
		"condition_else",
		"condition_then",
		"number_any_of",
		"number_one_of",
		"number_all_of",
		"number_not":
		return true
	default:
		return false
	}
}

// ValidateFile validates the contents of filePath with the schema
func ValidateFile(filePath string) (*gojsonschema.Result, error) {
	var collectionDir, err = ansible.GetStackHeadCollectionLocation()
	if err != nil {
		return nil, err
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + filepath.Join(collectionDir, "schema", "project-definition.schema.json"))

	configData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// YAML to JSON
	configJson, err := yaml.YAMLToJSON(configData)
	if err != nil {
		return nil, err
	}
	documentLoader := gojsonschema.NewBytesLoader(configJson)

	return gojsonschema.Validate(schemaLoader, documentLoader)
}

// ShouldValidate validates result
// signature uses interface{} and unused paremter because this method is also used in tests with Convey
func ShouldValidate(actual interface{}, _ ...interface{}) string {
	result := actual.(*gojsonschema.Result)
	if result.Valid() == true {
		return ""
	}
	errorMessage := fmt.Sprintf("The project definition is not valid. see errors:\n")

	for _, desc := range result.Errors() {
		if isInternalError(desc.Type()) {
			continue
		}
		errorMessage += fmt.Sprintf("- %s\n", desc)
	}
	return errorMessage
}
