package proxy_nginx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

type Data struct {
	ProjectName         string
	AllPortsTfString    string
	ServerConfig        string
	DependentContainers string
}

type Data2 struct {
	ProjectName string
	PortIndex   int
	Paths       Paths
	Domain      project.Domains
	Expose      project.DomainExpose
	Options     Options
}

type Paths struct {
	RootDirectory                string
	CertificatesDirectory        string
	CertificatesProjectDirectory string
	RootTerraformDirectory       string
	ProjectsRootDirectory        string
}

type Options struct {
	NginxUseHttps bool
}

type PortService struct {
	Expose                project.DomainExpose
	ContainerResourceName string
	Index                 int
}

func (p PortService) getTfString() string {
	return "${" + p.ContainerResourceName + ".ports[" + string(rune(p.Index)) + "].external}"
}

func buildSingleServerConfig(templateName string, portIndex int, expose project.DomainExpose, domain project.Domains) string {
	var files = []string{"./templates/nginx/nginx-base.conf.tmpl"}
	files = append(files, "./templates/nginx/nginx-"+templateName+".tmpl")

	paths := Paths{
		RootDirectory:                config.RootDirectory,
		CertificatesDirectory:        config.CertificatesDirectory,
		CertificatesProjectDirectory: system.Context.Project.GetCertificateDirectoryPath(),
		RootTerraformDirectory:       config.RootTerraformDirectory,
		ProjectsRootDirectory:        config.ProjectsRootDirectory,
	}

	var funcMap = template.FuncMap{
		"append": func(list []string, str string) []string {
			return append(list, str)
		},
		"dict_index_str": func(list []string, str string) int {
			for i, item := range list {
				if item == str {
					return i
				}
			}
			return -1
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, errors.New("invalid dictionary call")
			}

			root := make(map[string]interface{})

			for i := 0; i < len(values); i += 2 {
				dict := root
				var key string
				switch v := values[i].(type) {
				case string:
					key = v
				case []string:
					for i := 0; i < len(v)-1; i++ {
						key = v[i]
						var m map[string]interface{}
						v, found := dict[key]
						if found {
							m = v.(map[string]interface{})
						} else {
							m = make(map[string]interface{})
							dict[key] = m
						}
						dict = m
					}
					key = v[len(v)-1]
				default:
					return nil, errors.New("invalid dictionary key")
				}
				dict[key] = values[i+1]
			}

			return root, nil
		},
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
	}

	var tmpl = template.Must(template.New("").Funcs(funcMap).ParseFiles(files...))
	data := Data2{
		Paths:     paths,
		Domain:    domain,
		Expose:    expose,
		PortIndex: portIndex,
		Options: Options{
			NginxUseHttps: true,
		},
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		fmt.Println(err)
		return ""
	}

	return buf.String()
}

func buildServerConfig(project *project.Project, allServices []PortService) string {
	for _, domain := range project.Domains {
		for _, expose := range domain.Expose {
			if expose.ExternalPort == 443 {
				continue
			}
			portIndex := 0
			for _, service := range allServices {
				if expose.Service == service.Expose.Service {
					portIndex = service.Index
				}
			}

			if expose.ExternalPort == 80 {
				//{{ lookup('template', "{{ module_role_path | default(role_path) }}/templates/nginx/serverblock.http.j2", template_vars=dict(domainConfig=ns.domainCfg,expose=nginx_expose)) }}
				//{{ lookup('template', "{{ module_role_path | default(role_path) }}/templates/nginx/serverblock.https.j2", template_vars=dict(port_index=port_index,expose=nginx_expose,domainConfig=ns.domainCfg)) }}
				return buildSingleServerConfig("https", portIndex, expose, domain)
			} else {
				return buildSingleServerConfig("https", portIndex, expose, domain)
				//{{ lookup('template', "{{ module_role_path | default(role_path) }}/templates/nginx/serverblock.https.j2", template_vars=dict(port_index=port_index,expose=nginx_expose,port=nginx_expose.external_port,domainConfig=ns.domainCfg)) }}
			}
		}
	}
	return ""
}

func Deploy() {
	fmt.Println("Deploy step")
	var buf bytes.Buffer
	fileContents, err := os.ReadFile("./templates/nginx.tf.tmpl")
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("nginx").Parse(string(fileContents))
	if err != nil {
		panic(err)
	}

	var AllPorts []PortService
	var DependentContainers []string

	project := system.Context.Project
	for _, domain := range project.Domains {
		for i, expose := range domain.Expose {
			ContainerResourceName := "docker_container.stackhead-" + project.Name + "-" + expose.Service
			if expose.ExternalPort != 443 {
				DependentContainers = append(DependentContainers, ContainerResourceName)
			}
			//expose.Service
			AllPorts = append(AllPorts, PortService{
				Expose:                expose,
				ContainerResourceName: ContainerResourceName,
				Index:                 i,
			})
		}
	}

	var AllPortsTfStrings []string
	for _, port := range AllPorts {
		AllPortsTfStrings = append(AllPortsTfStrings, port.getTfString())
	}

	serverConfig := buildServerConfig(project, AllPorts)

	data := Data{
		ProjectName:         project.Name,
		AllPortsTfString:    strings.Join(AllPortsTfStrings, ","),
		ServerConfig:        serverConfig,
		DependentContainers: strings.Join(DependentContainers, ","),
	}
	err = tmpl.Execute(&buf, data)
}
