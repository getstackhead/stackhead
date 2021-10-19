package terraform

import (
	"path"

	"github.com/getstackhead/stackhead/config"
)

func GetCommand(command string) string {
	return "TF_DATA_DIR=" + path.Join(config.RootTerraformDirectory, ".terraform") + " terraform " + command
}
