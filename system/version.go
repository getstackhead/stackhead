package system

import (
	"fmt"
	"path"

	"github.com/blang/semver/v4"
	"github.com/getstackhead/stackhead/config"
	xfs "github.com/saitho/golang-extended-fs/v2"
)

var currentVersion = "2.0.0"
var remoteVersionFilePath = "ssh://" + path.Join(config.RootDirectory, "VERSION")

func WriteVersion() error {
	return xfs.WriteFile(remoteVersionFilePath, currentVersion)
}

func ValidateVersion() (bool, string, error) {
	remoteVersion, err := xfs.ReadFile(remoteVersionFilePath)
	if err != nil {
		return false, "", err
	}
	infoText := fmt.Sprintf("StackHead version used for setup is %s - Current version: %s", remoteVersion, currentVersion)
	//logger.Infoln(infoText)

	v1, err := semver.Make(remoteVersion)
	if err != nil {
		return false, infoText, err
	}
	v2, err := semver.Make(currentVersion)
	if err != nil {
		return false, infoText, err
	}

	return v1.Major == v2.Major, infoText, nil
}
