package commandsinit

import (
	"fmt"
	"strings"
	"sync"

	"github.com/getstackhead/stackhead/cli/routines"
	"github.com/getstackhead/stackhead/cli/stackhead"
)

func collectModules() []string {
	var modules []string
	var err error
	var moduleName string

	moduleName, err = stackhead.GetWebserverModule()
	if err != nil {
		panic(err.Error())
	}
	modules = append(modules, moduleName)

	moduleName, err = stackhead.GetContainerModule()
	if err != nil {
		panic(err.Error())
	}
	modules = append(modules, moduleName)

	pluginModules, err := stackhead.GetPluginModules()
	if err != nil {
		panic(err.Error())
	}
	modules = append(modules, pluginModules...)

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

// InstallModules is a list of task options that provide the actual workflow being run
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
