package jsonschema_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/getstackhead/stackhead/cli/jsonschema"
)

func collectYamlFiles(folder string) []string {
	var files []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func TestValidDefinitions(t *testing.T) {
	Convey("test valid definitions", t, func() {
		for _, file := range collectYamlFiles("../../validation/examples/valid") {
			Convey(fmt.Sprintf("file %s should validate", file), func() {
				result, err := jsonschema.ValidateFile(file)
				So(err, ShouldBeNil)
				So(result, jsonschema.ShouldValidate)
			})
		}
	})
}

func TestInvalidDefinitions(t *testing.T) {
	Convey("test invalid definitions", t, func() {
		for _, file := range collectYamlFiles("../../validation/examples/invalid") {
			Convey(fmt.Sprintf("file %s should not validate", file), func() {
				result, err := jsonschema.ValidateFile(file)
				So(err, ShouldBeNil)
				So(result, ShouldNotValidate)
			})
		}
	})
}

func ShouldNotValidate(actual interface{}, _ ...interface{}) string {
	result := jsonschema.ShouldValidate(actual)
	if result == "" { // file validated
		return "File validated (it should not)"
	} else {
		return ""
	}
}
