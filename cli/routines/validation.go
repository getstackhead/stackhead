package routines

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getstackhead/stackhead/cli/ansible"
	jsonschema "github.com/saitho/jsonschema-validator/validator"
)

func Validate(filePath string, schemaFile string)  {
	var collectionDir, err = ansible.GetStackHeadCollectionLocation()
	collectionAbsDir, err := filepath.Abs(collectionDir)
	if err != nil {
		panic(err)
		return
	}

	schemaPath := filepath.Join(collectionAbsDir, "schema", schemaFile)
	if err != nil {
		_, err = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		if err != nil {
			panic(err)
		}
		return
	}
	result, err := jsonschema.ValidateFile(filePath, schemaPath)

	if err != nil {
		panic(err.Error())
	}

	errorMessage := jsonschema.ShouldValidate(result)
	if len(errorMessage) == 0 {
		_, err = fmt.Fprintf(os.Stdout, "The project definition is valid.\n")
	} else {
		_, err = fmt.Fprintf(os.Stderr, errorMessage+"\n")
		if err != nil {
			panic(err.Error())
		}
		os.Exit(1)
	}
	if err != nil {
		panic(err.Error())
	}
}
