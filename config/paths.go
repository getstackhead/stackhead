package config

import (
	"net"
	"path"

	"github.com/shibukawa/configdir"
)

var RootDirectory = "/stackhead"
var CertificatesDirectory = RootDirectory + "/certificates"
var RootTerraformDirectory = RootDirectory + "/terraform"
var ProjectsRootDirectory = RootDirectory + "/projects"

func GetSnakeoilPaths() (string, string) {
	return path.Join(CertificatesDirectory, "fullchain_snakeoil.pem"), path.Join(CertificatesDirectory, "privkey_snakeoil.pem")
}

func GetLocalStackHeadConfigDir() string {
	configDirs := configdir.New("getstackhead", "stackhead")
	folders := configDirs.QueryFolders(configdir.Global)
	return folders[0].Path
}

func GetLocalRemoteKeyDir(host net.IP) string {
	return path.Join(GetLocalStackHeadConfigDir(), "ssh", "remotes", host.String())
}
