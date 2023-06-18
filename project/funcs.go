package project

import (
	"path"

	"github.com/getstackhead/stackhead/config"
)

func (project *Project) GetDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name)
}

func (project *Project) GetDeploymentsPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "deployments")
}

func (project *Project) GetRuntimeDataDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "data")
}
