package proxy

import (
	"strconv"
	"text/template"

	"github.com/getstackhead/stackhead/project"
)

type PortService struct {
	Expose                project.DomainExpose
	ContainerResourceName string
	Index                 int
}

func (p PortService) GetTfString() string {
	return "${" + p.ContainerResourceName + ".ports[" + strconv.Itoa(p.Index) + "].external}"
}

type ProxyRenderingContext struct {
	AllPorts []PortService
}

var Context = ProxyRenderingContext{}

// common functions used in proxy templates
var FuncMap = template.FuncMap{
	"getBasicAuths": func(s []project.DomainSecurityAuthentication) []project.DomainSecurityAuthentication {
		var auths []project.DomainSecurityAuthentication
		for _, authentication := range s {
			if authentication.Type != "basic" {
				continue
			}
			auths = append(auths, authentication)
		}
		return auths
	},
	"getPortIndex": func(service string, internalPort int) int {
		for _, port := range Context.AllPorts {
			if port.Expose.Service != service {
				continue
			}
			if port.Expose.InternalPort != internalPort {
				continue
			}
			return port.Index
		}
		return -1
	},
}
