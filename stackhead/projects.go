package stackhead

import (
	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
	xfs "github.com/saitho/golang-extended-fs"
)

type DeployedProject struct {
	Path    string
	Project project.Project
}

func GetDeployedProjects() ([]DeployedProject, error) {
	var deployedProjects []DeployedProject

	folders, err := xfs.ListFolders("ssh://" + config.ProjectsRootDirectory)
	if err != nil {
		if err.Error() == "file does not exist" {
			return deployedProjects, nil
		}
		return deployedProjects, err
	}
	for _, folder := range folders {
		println(folder.Name())
		deployedProjects = append(deployedProjects, DeployedProject{
			Path:    folder.Name(),
			Project: project.Project{}, // config.LoadProjectDefinition("ssh://" + folder.Name())
		})
	}

	return deployedProjects, nil
}
