package project

import (
	"path"

	"github.com/getstackhead/stackhead/config"
)

func (project *Project) GetDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name)
}

func (project *Project) GetCertificateDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "certificates")
}

func (project *Project) GetTerraformDirectoryPath() string {
	return path.Join(config.ProjectsRootDirectory, project.Name, "terraform")
}

func (project *Project) GetCertificatesDirectory() string {
	return path.Join(config.CertificatesDirectory, project.Name)
}