package container_docker

import (
	"fmt"

	docker_system "github.com/getstackhead/stackhead/modules/container/docker/system"
	"github.com/getstackhead/stackhead/system"
)

func (m Module) Destroy(modulesSettings interface{}) error {
	// Execute hooks
	if err := docker_system.ExecuteHook("beforeDestroy"); err != nil {
		return fmt.Errorf("Before destroy hook %s failed: ", err.Error())
	}

	// Stop and remove containers
	// todo: allow using either docker-compose or "docker compose" whichever is available (prefer "docker compose")
	_, stderr, err := system.RemoteRun("docker compose", system.RemoteRunOpts{Args: []string{"down"}, WorkingDir: system.Context.Project.GetDirectoryPath()})
	if err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("Unable to stop Docker containers: " + stderr.String())
		}
		return fmt.Errorf("Unable to stop Docker containers: " + err.Error())
	}
	return nil
}
