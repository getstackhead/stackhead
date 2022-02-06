package system

import (
	"fmt"
	"path"

	"github.com/blang/semver/v4"
	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/config"
)

var currentVersion = "2.0.0"
var remoteVersionFilePath = "ssh://" + path.Join(config.RootDirectory, "VERSION")

func WriteVersion() error {
	return xfs.WriteFile(remoteVersionFilePath, currentVersion)
}

func ValidateVersion() (bool, error) {
	remoteVersion, err := xfs.ReadFile(remoteVersionFilePath)
	if err != nil {
		return false, err
	}
	logger.Infoln(fmt.Sprintf("StackHead version used for setup is %s - Current version: %s", remoteVersion, currentVersion))

	v1, err := semver.Make(remoteVersion)
	if err != nil {
		return false, err
	}
	v2, err := semver.Make(currentVersion)
	if err != nil {
		return false, err
	}

	return v1.Major == v2.Major, nil
}
