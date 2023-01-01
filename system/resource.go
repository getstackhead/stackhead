package system

type Type string

const (
	TypeFile      Type = "file"
	TypeContainer Type = "container"
)

type Resource struct {
	Type Type
	
	// name of the resource (e.g. file name, container name)
	Name string

	// name of the associated service
	ServiceName string

	// for container resources
	Ports []string
}
