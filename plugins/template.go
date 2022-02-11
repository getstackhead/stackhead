package plugins

import (
	"bytes"
	"errors"
	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/system"
	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"
	"os"
	"path"
	"text/template"

	"github.com/getstackhead/stackhead/pluginlib"
)

func RenderTemplate(paths []string, data map[string]interface{}) (string, error) {
	var funcMap = template.FuncMap{
		"append": func(list []string, str string) []string {
			return append(list, str)
		},
		"dict_index_str": func(list []string, str string) int {
			for i, item := range list {
				if item == str {
					return i
				}
			}
			return -1
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dictionary call")
			}

			root := make(map[string]interface{})

			for i := 0; i < len(values); i += 2 {
				dict := root
				var key string
				switch v := values[i].(type) {
				case string:
					key = v
				case []string:
					for i := 0; i < len(v)-1; i++ {
						key = v[i]
						var m map[string]interface{}
						v, found := dict[key]
						if found {
							m = v.(map[string]interface{})
						} else {
							m = make(map[string]interface{})
							dict[key] = m
						}
						dict = m
					}
					key = v[len(v)-1]
				default:
					return nil, errors.New("invalid dictionary key")
				}
				dict[key] = values[i+1]
			}

			return root, nil
		},
		"getAuthsByType": func(authType string, s []pluginlib.DomainSecurityAuthentication) []pluginlib.DomainSecurityAuthentication {
			var auths []pluginlib.DomainSecurityAuthentication
			for _, authentication := range s {
				if authentication.Type != authType {
					continue
				}
				auths = append(auths, authentication)
			}
			return auths
		},
	}

	fileContent := ""
	for _, templatePath := range paths {
		file, err := os.ReadFile(templatePath)
		if err != nil {
			return "", err
		}
		fileContent += string(file)
	}

	var tmpl = template.Must(template.New("tpl").Funcs(funcMap).Parse(fileContent))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", errors.New("error processing the plugin template: " + err.Error())
	}

	return buf.String(), nil
}

func CreateTerraformFile(fileName string, fileContent string) error {
	err := xfs.WriteFile("ssh://"+path.Join(
		config.Paths.GetProjectTerraformDirectoryPath(system.Context.Project),
		fileName,
	), fileContent)
	logger.Info("Writing files to")
	return err
}
