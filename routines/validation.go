package routines

import (
	"crypto/tls"
	"fmt"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	jsonschema "github.com/saitho/jsonschema-validator/validator"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"net/http"
	"os"
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
		// Use schema stored in binary
		var f pkging.File
		f, err = pkger.Open("/schemas/" + schemaFile)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()

		var sl []byte
		sl, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err.Error())
		}
		result, err = jsonschema.ValidateFileWithInput(filePath, sl)
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
