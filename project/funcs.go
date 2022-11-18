package project

import (
	"path"

	"github.com/getstackhead/stackhead/config"
)

func (project *Project) GetDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name)
}

func (project *Project) GetRuntimeDataDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "runtime")
}
