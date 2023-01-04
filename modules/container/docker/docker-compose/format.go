package docker_compose

import (
	"fmt"
	"regexp"

	"github.com/getstackhead/stackhead/project"
)

func GetVolumeSrcKey(projectName string, serviceName string, volume project.ContainerServiceVolume) string {
	name := projectName
	if volume.Type == "local" {
		name += "-" + serviceName
	}
	return fmt.Sprintf("%s-%s-%s", volume.Type, name, SanitizeVolume(volume.Src))
}

func SanitizeVolume(s string) string {
	var re = regexp.MustCompile(`[^\w]`)
	return re.ReplaceAllString(s, "_")
}
