package system

import (
	"fmt"
	"path"
	"regexp"
	"sort"

	xfs "github.com/saitho/golang-extended-fs/v2"
	"gopkg.in/yaml.v3"
	"path/filepath"

	"github.com/getstackhead/stackhead/project"
)

func GetLatestDeployment(project *project.Project) (*Deployment, error) {
	files, err := xfs.ListFolders("ssh://" + project.GetDeploymentsPath())
	if err != nil {
		return nil, err
	}
	if files != nil {
		// newest files at the top
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().After(files[j].ModTime())
		})
		for _, file := range files {
			if file.IsDir() && MatchDeploymentNaming(file.Name()) {
				fullPath := path.Join(project.GetDeploymentsPath(), file.Name())
				latestDeployment, err := GetDeploymentByPath(fullPath)
				if err != nil {
					return nil, err
				}
				if !latestDeployment.RolledBack {
					latestDeployment.Project = project
					return latestDeployment, nil
				}
			}
		}
	}
	return nil, nil
}

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
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal([]byte(deploymentFile), &deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
}
