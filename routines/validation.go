package routines

import (
	"crypto/tls"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/markbates/pkger"
	xfs "github.com/saitho/golang-extended-fs"
	jsonschema "github.com/saitho/jsonschema-validator/validator"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

type ValidationSource string

func CobraValidationBase(schemaFile string, version string, branch string, ignoreSslCertificate bool) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		source := "stackhead_cli"
		if len(version) > 0 {
			source = "https://schema.stackhead.io/stackhead-cli/tag/" + version + "/-"
		} else if len(branch) > 0 {
			source = "https://schema.stackhead.io/stackhead-cli/branch/" + branch + "/-"
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			MinVersion:         1,
			InsecureSkipVerify: ignoreSslCertificate, // nolint:gosec
		}
		Validate(args[0], schemaFile, source)
	}
}

func Validate(filePath string, schemaFile string, source string) {
	var err error
	var result *gojsonschema.Result

	switch source {
	case "stackhead_cli":
		var tempDirName string
		// Use schema stored in binary
		tempDirName, err = os.MkdirTemp("", "")
		if err != nil {
			panic("unable to create temporary folder: " + err.Error())
		}
		defer os.RemoveAll(tempDirName)

		// Workaround to resolve references correctly: copy everything to local file system
		err = pkger.Walk("/schemas/", func(filePath string, info fs.FileInfo, err error) error {
			actualFilePath := strings.TrimPrefix(filePath, "github.com/getstackhead/stackhead:")
			fullFilePath := path.Join(tempDirName, actualFilePath)
			if info.IsDir() {
				if err := os.MkdirAll(fullFilePath, os.ModeDir|os.ModePerm); err != nil {
					return fmt.Errorf("unable to create temporary folder at " + fullFilePath)
				}
			} else {
				if err := xfs.CopyFile("pkging://"+actualFilePath, path.Join(fullFilePath)); err != nil {
					return fmt.Errorf("unable to create temporary file at " + fullFilePath + ": " + err.Error())
				}
			}
			return nil
		})

		if err != nil {
			panic(err.Error())
		}
		result, err = jsonschema.ValidateFile(filePath, path.Join(tempDirName, "schemas", schemaFile))
		if err != nil {
			panic(err.Error())
		}
	default:
		url := fmt.Sprintf("%s/%s", source, schemaFile)
		fmt.Fprintf(os.Stdout, "Validating with schema from URL: %s\n", url)
		// Pull from online Schemastore, source contains the URL
		result, err = jsonschema.ValidateFile(filePath, url)
		if err != nil {
			panic(err.Error())
		}
	}

	errorMessage := jsonschema.ShouldValidate(result)
	if len(errorMessage) == 0 {
		_, err = fmt.Fprintf(os.Stdout, "The project definition is valid.\n")
	} else {
		_, err = fmt.Fprintf(os.Stderr, errorMessage+"\n")
		if err != nil {
			panic(err.Error())
		}
		defer func() { os.Exit(1) }()
	}
	if err != nil {
		panic(err.Error())
	}
}
