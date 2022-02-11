package config

import (
	"net"
	"path"

	"github.com/shibukawa/configdir"

	"github.com/getstackhead/stackhead/pluginlib"
)

type PathsStruct struct {
	Root          string
	Certificates  string
	RootTerraform string
	ProjectsRoot  string
}

var Paths = PathsStruct{
	Root:          "/stackhead",
	Certificates:  "/stackhead/certificates",
	RootTerraform: "/stackhead/terraform",
	ProjectsRoot:  "/stackhead/projects",
}

func (p PathsStruct) GetSnakeoilFullchainPath() string {
	return path.Join(p.Certificates, "fullchain_snakeoil.pem")
}

func (p PathsStruct) GetSnakeoilPrivKeyPath() string {
	return path.Join(p.Certificates, "privkey_snakeoil.pem")
}

func (p PathsStruct) GetProjectCertificateDirectoryPath(project *pluginlib.Project) string {
	return path.Join(p.Certificates, project.Name)
}

func (p PathsStruct) GetProjectDirectoryPath(project *pluginlib.Project) string {
	return path.Join(p.ProjectsRoot, project.Name)
}

func (p PathsStruct) GetProjectTerraformDirectoryPath(project *pluginlib.Project) string {
	return path.Join(p.ProjectsRoot, project.Name, "terraform")
}

func (p PathsStruct) GetCertificatesDirectoryForProject(project *pluginlib.Project) string {
	return path.Join(p.Certificates, project.Name)
}

func GetLocalStackHeadConfigDir() string {
	configDirs := configdir.New("getstackhead", "stackhead")
	folders := configDirs.QueryFolders(configdir.Global)
	return folders[0].Path
}

func GetLocalRemoteKeyDir(host net.IP) string {
	return path.Join(GetLocalStackHeadConfigDir(), "ssh", "remotes", host.String())
}
