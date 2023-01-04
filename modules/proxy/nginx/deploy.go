package proxy_nginx

import (
	"fmt"
	xfs "github.com/saitho/golang-extended-fs/v2"
	"path"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/modules/proxy"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

type Paths struct {
	RootDirectory                string
	CertificatesProjectDirectory string
	ProjectCertificatesDirectory string
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
		ProjectsRootDirectory:        config.ProjectsRootDirectory,
		AcmeChallengesDirectory:      AcmeChallengesDirectory,
		SnakeoilFullchainPath:        SnakeoilFullchainPath,
		SnakeoilPrivkeyPath:          SnakeoilPrivkeyPath,
	}
}

func buildSingleServerConfig(templateName string, portIndex int, expose project.DomainExpose, domain project.Domains) string {
	var files = []string{"nginx/nginx-base.conf.tmpl"}
	files = append(files, "nginx/nginx-"+templateName+".tmpl")

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
		fsContent, err := templates.ReadFile("templates/" + file)
		if err != nil {
			panic("Unable to read template file \"" + file + "\": " + err.Error())
		}
		text += string(fsContent)
	}

	data := map[string]any{
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
		panic("Unable to build Nginx proxy template: " + err.Error())
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
	moduleSettings.Config.SetDefaults()

	fmt.Println("Deploy step")
	paths := getPaths()

	if err := xfs.CreateFolder("ssh://" + paths.ProjectCertificatesDirectory); err != nil {
		return err
	}
	if err := xfs.CreateFolder("ssh://" + paths.CertificatesProjectDirectory); err != nil {
		return err
	}

	serverConfig := buildServerConfig(system.Context.Project, proxy.Context.AllPorts)
	nginxProjectLocation := system.Context.Project.GetDirectoryPath() + "/nginx.conf"
	if err = xfs.WriteFile("ssh://"+nginxProjectLocation, serverConfig); err != nil {
		return err
	}

	// Symlink project certificate files to snakeoil files after initial creation
	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-s " + paths.SnakeoilFullchainPath + " " + paths.CertificatesProjectDirectory + "/fullchain.pem"}, AllowFail: true}); err != nil {
		return fmt.Errorf("Unable to symlink snakeoil full chain: " + err.Error())
	}
	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-s " + paths.SnakeoilPrivkeyPath + " " + paths.CertificatesProjectDirectory + "/privkey.pem"}, AllowFail: true}); err != nil {
		return fmt.Errorf("Unable to symlink snakeoil privkey: " + err.Error())
	}

	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf " + nginxProjectLocation + " /etc/nginx/sites-available/stackhead_" + system.Context.Project.Name + ".conf"}}); err != nil {
		return fmt.Errorf("Unable to symlink project Nginx file: " + err.Error())
	}
	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf /etc/nginx/sites-available/stackhead_" + system.Context.Project.Name + ".conf " + moduleSettings.Config.VhostPath + "/stackhead_" + system.Context.Project.Name + ".conf"}}); err != nil {
		return fmt.Errorf("Unable to enable project Nginx file: " + err.Error())
	}
	// first reload so webserver config works for ACME request
	if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "nginx"}, Sudo: true}); err != nil {
		return fmt.Errorf("Unable to reload Nginx service: " + err.Error())
	}

	certMail := "certificates-noreply@stackhead.io"
	if len(moduleSettings.CertificatesEmail) > 0 {
		certMail = moduleSettings.CertificatesEmail
	}
	if err := generateCertificates(paths, certMail); err != nil {
		return fmt.Errorf("Unable to generate certificates: " + err.Error())
	}
	// reload Nginx again so certificates take effect
	if _, err := system.SimpleRemoteRun("systemctl", system.RemoteRunOpts{Args: []string{"reload", "nginx"}, Sudo: true}); err != nil {
		return fmt.Errorf("Unable to reload Nginx service: " + err.Error())
	}

	return nil
}

/**
 * Create certificate files and remove Nginx configuration for ACME confirmation
 */
func generateCertificates(paths Paths, certMail string) error {
	// Create AcmeChallengesDirectory folder
	firstDomain := system.Context.Project.Domains[0].Domain
	domainChallengeDir := path.Join(AcmeChallengesDirectory, firstDomain)
	if err := xfs.CreateFolder("ssh://" + domainChallengeDir); err != nil {
		return err
	}
	if err := xfs.Chown("ssh://"+AcmeChallengesDirectory, 1412, 1412); err != nil {
		return err
	}

	args := []string{
		"certonly",
		"-m " + certMail,
		"--no-eff-email",
		"--agree-tos",
		"-q",
		"--webroot",
		"-w " + domainChallengeDir,
	}
	for _, domain := range system.Context.Project.Domains {
		args = append(args, "-d "+domain.Domain)
	}
	if system.Context.IsCI {
		args = append(args, "--test-cert")
	}

	if result, err := system.SimpleRemoteRun("certbot", system.RemoteRunOpts{Args: args, Sudo: true}); err != nil {
		fmt.Println(result)
		return err
	}

	// Overwrite symlinked snakeoil certificates
	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf /etc/letsencrypt/live/" + firstDomain + "/fullchain.pem " + paths.CertificatesProjectDirectory + "/fullchain.pem"}}); err != nil {
		return fmt.Errorf("Unable to symlink snakeoil full chain: " + err.Error())
	}
	if _, err := system.SimpleRemoteRun("ln", system.RemoteRunOpts{Args: []string{"-sf /etc/letsencrypt/live/" + firstDomain + "/privkey.pem " + paths.CertificatesProjectDirectory + "/privkey.pem"}}); err != nil {
		return fmt.Errorf("Unable to symlink snakeoil privkey: " + err.Error())
	}

	return nil
}
