package config

import (
	"net"
	"path"

	"github.com/shibukawa/configdir"
)

var RootDirectory = "/etc/stackhead"
var ProjectsRootDirectory = RootDirectory + "/projects"

func GetServerConfigFilePath() string {
	return path.Join(RootDirectory, "config.yml")
}

func GetLocalStackHeadConfigDir() string {
	configDirs := configdir.New("getstackhead", "stackhead")
	folders := configDirs.QueryFolders(configdir.Global)
	return folders[0].Path
}

func GetLocalRemoteKeyDir(host net.IP) string {
	return path.Join(GetLocalStackHeadConfigDir(), "ssh", "remotes", host.String())
}
