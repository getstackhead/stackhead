package config

import (
	"net"
	"path"

	"github.com/shibukawa/configdir"

	"github.com/getstackhead/stackhead/pluginlib"
)

var RootDirectory = "/stackhead"
var CertificatesDirectory = RootDirectory + "/certificates"
var RootTerraformDirectory = RootDirectory + "/terraform"
var ProjectsRootDirectory = RootDirectory + "/projects"

func GetLocalStackHeadConfigDir() string {
	configDirs := configdir.New("getstackhead", "stackhead")
	folders := configDirs.QueryFolders(configdir.Global)
	return folders[0].Path
}

func GetLocalRemoteKeyDir(host net.IP) string {
	return path.Join(GetLocalStackHeadConfigDir(), "ssh", "remotes", host.String())
}

func GetProjectDirectoryPath(project *pluginlib.Project) string {
	return path.Join(ProjectsRootDirectory, project.Name)
}

func GetProjectCertificateDirectoryPath(project *pluginlib.Project) string {
	return path.Join(ProjectsRootDirectory, project.Name, "certificates")
}

func GetProjectTerraformDirectoryPath(project *pluginlib.Project) string {
	return path.Join(ProjectsRootDirectory, project.Name, "terraform")
}

func GetCertificatesDirectoryForProject(project *pluginlib.Project) string {
	return path.Join(CertificatesDirectory, project.Name)
}
