package system

import (
	"fmt"
	"regexp"

	xfs "github.com/saitho/golang-extended-fs/v2"
	"gopkg.in/yaml.v3"
	"path/filepath"
)

func MatchDeploymentNaming(folderName string) bool {
	pattern := regexp.MustCompile(`(?m)^v(\d+)$`)
	return pattern.MatchString(folderName)
}

func GetDeploymentByPath(path string) (*Deployment, error) {
	deployment := Deployment{}
	if !MatchDeploymentNaming(filepath.Base(path)) {
		return nil, fmt.Errorf("last folder in path should be a version folder")
	}
	deploymentFilePath := "ssh://" + filepath.Join(path, "deployment.yaml")
	hasDeploymentFile, err := xfs.HasFile(deploymentFilePath)
	if err != nil {
		return nil, err
	}
	if !hasDeploymentFile {
		return nil, fmt.Errorf("Missing deployment file in folder")
	}
	deploymentFile, err := xfs.ReadFile(deploymentFilePath)
	if err = yaml.Unmarshal([]byte(deploymentFile), &deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
}
