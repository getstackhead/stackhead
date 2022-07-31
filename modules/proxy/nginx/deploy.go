package proxy_nginx

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"strings"
	"text/template"

	xfs "github.com/saitho/golang-extended-fs"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/modules/proxy"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

type Data struct {
	ProjectName         string
	AllPortsTfString    string
	ServerConfig        string
	DependentContainers string
	Paths               Paths
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
	SnakeoilFullchainPath        string
	SnakeoilPrivkeyPath          string
}

type Options struct {
	NginxUseHttps bool
}

func getPaths() Paths {
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := config.GetSnakeoilPaths()
	return Paths{
		RootDirectory:                config.RootDirectory,
		CertificatesDirectory:        config.CertificatesDirectory,
		CertificatesProjectDirectory: system.Context.Project.GetCertificateDirectoryPath(),
		RootTerraformDirectory:       config.RootTerraformDirectory,
		ProjectsRootDirectory:        config.ProjectsRootDirectory,
		SnakeoilFullchainPath:        SnakeoilFullchainPath,
		SnakeoilPrivkeyPath:          SnakeoilPrivkeyPath,
	}
}

func buildSingleServerConfig(templateName string, portIndex int, expose project.DomainExpose, domain project.Domains) string {
	var files = []string{"pkging:///templates/modules/proxy/nginx/nginx/nginx-base.conf.tmpl"}
	files = append(files, "pkging:///templates/modules/proxy/nginx/nginx/nginx-"+templateName+".tmpl")

	var funcMap = template.FuncMap{
		"dict_index_str": func(list []string, str string) int {
			for i, item := range list {
				if item == str {
					return i
				}
			}
			return -1
		},
	}

	var text string
	for _, file := range files {
		fsContent, err := xfs.ReadFile(file)
		if err != nil {
			panic("Unable to read template file \"" + file + "\": " + err.Error())
		}
		text += fsContent
	}

	var tmpl = template.Must(template.New("").Funcs(funcMap).Funcs(proxy.FuncMap).Funcs(sprig.TxtFuncMap()).Parse(text))
	data := Data2{
		Paths:     getPaths(),
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

func buildServerConfig(project *project.Project, allServices []proxy.PortService) string {
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
				httpConfig := buildSingleServerConfig("http", portIndex, expose, domain)
				expose.ExternalPort = 443
				httpsConfig := buildSingleServerConfig("https", portIndex, expose, domain)
				return httpConfig + "\n" + httpsConfig
			} else {
				return buildSingleServerConfig("https", portIndex, expose, domain)
				//{{ lookup('template', "{{ module_role_path | default(role_path) }}/templates/nginx/serverblock.https.j2", template_vars=dict(port_index=port_index,expose=nginx_expose,port=nginx_expose.external_port,domainConfig=ns.domainCfg)) }}
			}
		}
	}
	return ""
}

func (NginxProxyModule) Deploy() error {
	fmt.Println("Deploy step")
	var buf bytes.Buffer
	fileContents, err := xfs.ReadFile("pkging:///templates/modules/proxy/nginx/nginx.tf.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("nginx").Parse(fileContents)
	if err != nil {
		return err
	}

	var DependentContainers []string

	for _, domain := range system.Context.Project.Domains {
		for i, expose := range domain.Expose {
			ContainerResourceName := "docker_container.stackhead-" + system.Context.Project.Name + "-" + expose.Service
			if expose.ExternalPort != 443 {
				DependentContainers = append(DependentContainers, ContainerResourceName)
			}
			//expose.Service
			proxy.Context.AllPorts = append(proxy.Context.AllPorts, proxy.PortService{
				Expose:                expose,
				ContainerResourceName: ContainerResourceName,
				Index:                 i,
			})
		}
	}

	var AllPortsTfStrings []string
	for _, port := range proxy.Context.AllPorts {
		AllPortsTfStrings = append(AllPortsTfStrings, port.GetTfString())
	}

	serverConfig := buildServerConfig(system.Context.Project, proxy.Context.AllPorts)

	data := Data{
		ProjectName:         system.Context.Project.Name,
		AllPortsTfString:    strings.Join(AllPortsTfStrings, ","),
		ServerConfig:        serverConfig,
		DependentContainers: strings.Join(DependentContainers, ","),
		Paths:               getPaths(),
	}
	err = tmpl.Execute(&buf, data)

	if err != nil {
		return err
	}

	err = xfs.WriteFile("ssh://"+system.Context.Project.GetTerraformDirectoryPath()+"/nginx.tf", buf.String())

	// Todo: generate ssl certificates

	return err
}
