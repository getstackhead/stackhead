package proxy_nginx

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	xfs "github.com/saitho/golang-extended-fs"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/modules/proxy"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

type Paths struct {
	RootDirectory                string
	CertificatesProjectDirectory string
	ProjectCertificatesDirectory string
	RootTerraformDirectory       string
	ProjectsRootDirectory        string
	AcmeChallengesDirectory      string
	SnakeoilFullchainPath        string
	SnakeoilPrivkeyPath          string
}

type Options struct {
	NginxUseHttps bool
}

func getPaths() Paths {
	SnakeoilFullchainPath, SnakeoilPrivkeyPath := GetSnakeoilPaths()
	return Paths{
		RootDirectory:                config.RootDirectory,
		CertificatesProjectDirectory: GetCertificatesDirectory(system.Context.Project),
		ProjectCertificatesDirectory: GetCertificateDirectoryPath(system.Context.Project),
		RootTerraformDirectory:       config.RootTerraformDirectory,
		ProjectsRootDirectory:        config.ProjectsRootDirectory,
		AcmeChallengesDirectory:      AcmeChallengesDirectory,
		SnakeoilFullchainPath:        SnakeoilFullchainPath,
		SnakeoilPrivkeyPath:          SnakeoilPrivkeyPath,
	}
}

func buildSingleServerConfig(templateName string, portIndex int, expose project.DomainExpose, domain project.Domains) string {
	var files = []string{"pkging:///templates/modules/proxy/nginx/nginx/nginx-base.conf.tmpl"}
	files = append(files, "pkging:///templates/modules/proxy/nginx/nginx/nginx-"+templateName+".tmpl")

	var funcMap = proxy.FuncMap
	funcMap["dict_index_str"] = func(list []string, str string) int {
		for i, item := range list {
			if item == str {
				return i
			}
		}
		return -1
	}

	var text string
	for _, file := range files {
		fsContent, err := xfs.ReadFile(file)
		if err != nil {
			panic("Unable to read template file \"" + file + "\": " + err.Error())
		}
		text += fsContent
	}

	data := map[string]interface{}{
		"Paths":     getPaths(),
		"Domain":    domain,
		"Expose":    expose,
		"PortIndex": portIndex,
		"Options": Options{
			NginxUseHttps: true,
		},
	}
	nginxServerBlock, err := system.RenderModuleTemplateText(
		"base",
		text,
		data,
		funcMap)
	if err != nil {
		return ""
	}
	return nginxServerBlock
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
			}
		}
	}
	return ""
}

func (Module) Deploy(_modulesSettings interface{}) error {
	moduleSettings, err := system.UnpackModuleSettings[ModuleSettings](_modulesSettings)
	if err != nil {
		return fmt.Errorf("unable to load module settings: " + err.Error())
	}
	fmt.Println("Deploy step")
	paths := getPaths()

	if err := xfs.CreateFolder("ssh://" + paths.ProjectCertificatesDirectory); err != nil {
		return err
	}
	if err := xfs.CreateFolder("ssh://" + paths.CertificatesProjectDirectory); err != nil {
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

	data := map[string]interface{}{
		"AllPortsTfString":        strings.Join(AllPortsTfStrings, ","),
		"ServerConfig":            serverConfig,
		"DependentContainers":     strings.Join(DependentContainers, ","),
		"Paths":                   paths,
		"CertificatesMailAddress": moduleSettings.CertificatesEmail,
	}

	nginxTemplate, err := system.RenderModuleTemplate(
		"proxy/nginx/nginx.tf.tmpl",
		data,
		nil)
	if err != nil {
		return err
	}
	err = xfs.WriteFile("ssh://"+system.Context.Project.GetTerraformDirectoryPath()+"/nginx.tf", nginxTemplate)

	return generateCertificates(data)
}

/**
 * Create certificate files and remove Nginx configuration for ACME confirmation
 */
func generateCertificates(data map[string]interface{}) error {
	funcMap := template.FuncMap{
		"GetDomainNames": func(domains []project.Domains, start int) string {
			var names []string
			for i, domain := range domains {
				if i < start {
					continue
				}
				names = append(names, domain.Domain)
			}
			output, _ := json.Marshal(names)
			return string(output)
		},
	}
	acmeResolver, err := system.RenderModuleTemplate(
		"proxy/nginx/certificates/acme_challenge_resolver.sh.tmpl",
		data,
		funcMap)
	if err != nil {
		return err
	}
	resolverRemotePath := "ssh://" + system.Context.Project.GetDirectoryPath() + "/acme_challenge_resolver.sh"
	if err := xfs.WriteFile(resolverRemotePath, acmeResolver); err != nil {
		return err
	}
	if err := xfs.Chmod(resolverRemotePath, 0775); err != nil {
		return err
	}

	sslCertFile, err := system.RenderModuleTemplate("proxy/nginx/certificates/ssl-certificate.tf.tmpl", data, funcMap)
	if err != nil {
		return err
	}
	return xfs.WriteFile("ssh://"+system.Context.Project.GetTerraformDirectoryPath()+"/ssl-certificate.tf", sslCertFile)
}
