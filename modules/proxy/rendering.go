package proxy

import (
	"strconv"
	"strings"
	"text/template"

	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
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
	"getExternalPort": func(service string, internalPort int) string {
		for _, resource := range system.Context.Resources {
			if resource.Type != system.TypeContainer {
				continue
			}
			if resource.ServiceName != service {
				continue
			}
			for _, port := range resource.Ports {
				split := strings.SplitN(port, ":", 2)
				if split[1] == strconv.Itoa(internalPort) {
					return split[0]
				}
			}
		}
		return ""
	},
}
