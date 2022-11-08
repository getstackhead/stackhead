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

func (project *Project) GetTerraformDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "terraform")
}

func (project *Project) GetTerraformProjectProvidersFilePath() string {
	return path.Join(project.GetTerraformDirectoryPath(), "terraform-providers-project.tf")
}
