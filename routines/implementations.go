package routines

import (
	"fmt"
	"github.com/getstackhead/stackhead/project"
	xfs "github.com/saitho/golang-extended-fs/v2"

	logger "github.com/sirupsen/logrus"

	"github.com/getstackhead/stackhead/system"
)

var ValidateStackHeadVersionTask = Task{
	Name: fmt.Sprintf("Validating StackHead version"),
	Run: func(r RunningTask) error {
		isValid, infoText, err := system.ValidateVersion()
		r.SetSuccessMessage(infoText)
		if err != nil {
			logger.Debugln(err)
			err = fmt.Errorf("unable to validate StackHead version.")
		}
		if !isValid {
			err = fmt.Errorf("Trying to deploy with a newer version of StackHead than used for server setup. Please run a server setup again.")
		}

		if err != nil {
			r.SetFailMessage(err.Error())
		}
		return err
	},
}

var PrepareProjectTask = func(projectDefinition *project.Project) Task {
	return Task{
		Name: fmt.Sprintf("Preparing project structure"),
		Run: func(r RunningTask) error {
			r.PrintLn("Create project directory if not exists")
			if err := xfs.CreateFolder("ssh://" + projectDefinition.GetDirectoryPath()); err != nil {
				return err
			}
			return nil
		},
		ErrorAsErrorMessage: true,
	}
}
var CollectResourcesTask = func(projectDefinition *project.Project) Task {
	return Task{
		Name: fmt.Sprintf("Collecting resources"),
		Run: func(r RunningTask) error {
			// Collect exposed services
			var exposedServices []project.DomainExpose
			for _, domain := range projectDefinition.Domains {
				exposedServices = append(exposedServices, domain.Expose...)
			}
			for _, module := range system.Context.GetModulesInOrder() {
				if module.GetConfig().Type == "plugin" {
					continue
				}
				moduleSettings := system.GetModuleSettings(module.GetConfig().Name)
				if err := module.Deploy(moduleSettings); err != nil {
					return err
				}
			}
			return nil
		},
	}
}
