package system

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
	Ports []string

	ApplyResourceFunc    ApplyResourceFuncType
	RollbackResourceFunc RollbackResourceFuncType
}
