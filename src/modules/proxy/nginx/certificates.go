package proxy_nginx

import (
	"path"

	"github.com/getstackhead/stackhead/config"
	"github.com/getstackhead/stackhead/project"
)

var CertificatesDirectory = config.RootDirectory + "/certificates"
var AcmeChallengesDirectory = config.RootDirectory + "/acme-challenges"

func GetSnakeoilPaths() (string, string) {
	return path.Join(CertificatesDirectory, "fullchain_snakeoil.pem"), path.Join(CertificatesDirectory, "privkey_snakeoil.pem")
}

func GetCertificateDirectoryPath(p *project.Project) string {
	return path.Join(config.ProjectsRootDirectory, p.Name, "certificates")
}

func GetCertificatesDirectory(p *project.Project) string {
	return path.Join(CertificatesDirectory, p.Name)
}
