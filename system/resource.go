package system

import (
	"fmt"
	"strings"

	xfs "github.com/saitho/golang-extended-fs/v2"
)

type Type string

const (
	TypeFile      Type = "file"
	TypeFolder    Type = "folder"
	TypeLink      Type = "link"
	TypeContainer Type = "container"
)

type Operation string

const (
	OperationCreate Operation = "create"
	OperationDelete Operation = "delete"
)

type ApplyResourceFuncType func() error
type RollbackResourceFuncType func() error

type ResourceGroup struct {
	Name      string
	Resources []Resource

	ApplyResourceFunc    ApplyResourceFuncType    `yaml:"-"`
	RollbackResourceFunc RollbackResourceFuncType `yaml:"-"`
}

type Resource struct {
	Type           Type
	Operation      Operation `yaml:"-"`
	BackupFilePath string    `yaml:"-"`

	// if set the Name refers to an external resource. for files an absolute path is expected
	ExternalResource bool `yaml:"externalResource,omitempty"`

	// name of the resource (e.g. file name, container name)
	Name string

	// contents of resource, if any
	Content string `yaml:"-"`

	// name of the associated service
	ServiceName string `yaml:"serviceName,omitempty"`

	// for container resources
	Ports     []string `yaml:"ports,omitempty"`
	ImageName string   `yaml:"imageName,omitempty"`

	// for link resources
	LinkSource string `yaml:"linkSource,omitempty"`
	// EnforceLink is true, then symlink is created with force and will not ignore errors
	EnforceLink bool `yaml:"-"`
}

func (r Resource) GetOperationLabel(invertOperation bool) string {
	operation := "UNKNOWN"
	if r.Operation == OperationCreate {
		operation = "CREATE"
		if r.Type == TypeFile {
			// check if file exists
			hasFile, _ := xfs.HasFile("ssh://" + r.Name)
			if hasFile {
				operation = "UPDATE"
			}
		}
	} else if r.Operation == OperationDelete {
		operation = "DELETE"
	}

	if invertOperation {
		if operation == "CREATE" {
			operation = "DELETE"
		}
	}

	return operation
}

func (r Resource) ToString(invertOperation bool) string {
	if r.Type == TypeContainer {
		return fmt.Sprintf("[%s] Container %s (service=%s, image=%s, ports=%s)", r.GetOperationLabel(invertOperation), r.Name, r.ServiceName, r.ImageName, strings.Join(r.Ports, ", "))
	}
	return fmt.Sprintf("[%s] %s %s", r.GetOperationLabel(invertOperation), r.Type, r.Name)
}
