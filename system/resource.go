package system

import (
	"fmt"
	"strings"

	xfs "github.com/saitho/golang-extended-fs/v2"
)

type Type string

const (
	TypeFile      Type = "file"
	TypeContainer Type = "container"
)

type Operation string

const (
	OperationCreate Operation = "create"
)

type ApplyResourceFuncType func() error
type RollbackResourceFuncType func() error

type ResourceGroup struct {
	Name      string
	Resources []Resource

	ApplyResourceFunc    ApplyResourceFuncType
	RollbackResourceFunc RollbackResourceFuncType
}

type Resource struct {
	Type      Type
	Operation Operation

	// name of the resource (e.g. file name, container name)
	Name string

	// contents of resource, if any
	Content string

	// name of the associated service
	ServiceName string

	// for container resources
	Ports     []string
	ImageName string
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
