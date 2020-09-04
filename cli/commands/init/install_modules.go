package commands_init

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/getstackhead/stackhead/cli/routines"
	"github.com/getstackhead/stackhead/cli/stackhead"
)

func collectModules() []string {
	var modules []string
	var err error

	var webserver = viper.GetString("modules.webserver")
	if len(webserver) == 0 {
		webserver = "getstackhead.stackhead_webserver_nginx"
	}
	webserver, err = stackhead.AutoCompleteModuleName(webserver, stackhead.ModuleWebserver)
	if err != nil {
		// error
	}
	modules = append(modules, webserver)

	var container = viper.GetString("modules.container")
	if len(container) == 0 {
		container = "getstackhead.stackhead_container_docker"
	}
	container, err = stackhead.AutoCompleteModuleName(container, stackhead.ModuleContainer)
	if err != nil {
		// error
	}

	modules = append(modules, container)
	return modules
}

func installStackHeadModules() error {
	var modules = collectModules()

	var wg sync.WaitGroup
	wg.Add(len(modules))
	allErrorsChan := make(chan routines.TaskResult, len(modules))
	for _, moduleName := range modules {
		go func(name string) {
			defer wg.Done()
			err := routines.ExecAnsibleGalaxy("install", name)
			if err != nil {
				allErrorsChan <- routines.TaskResult{
					Name:    name,
					Message: err.Error(),
					Error:   true,
				}
			}
		}(moduleName)
	}
	wg.Wait()
	close(allErrorsChan)

	if len(allErrorsChan) > 0 {
		var readableErrors []string
		for err := range allErrorsChan {
			readableErrors = append(readableErrors, fmt.Sprintf("- %s: %s ", err.Name, err.Message))
		}
		return fmt.Errorf("The following errors occurred while installing StackHead modules:\n%s", strings.Join(readableErrors, "\n"))
	}
	return nil
}

var InstallModules = []routines.TaskOption{
	routines.Text("Installing StackHead modules"),
	routines.Execute(func(wg *sync.WaitGroup, result chan routines.TaskResult) {
		defer wg.Done()

		err := installStackHeadModules()

		taskResult := routines.TaskResult{
			Error: err != nil,
		}
		if err != nil {
			taskResult.Message = err.Error()
		}

		result <- taskResult
	}),
}
