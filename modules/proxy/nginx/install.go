package proxy_nginx

import (
	"fmt"
	"github.com/getstackhead/stackhead/system"

	xfs "github.com/saitho/golang-extended-fs"
	logger "github.com/sirupsen/logrus"
)

func (NginxProxyModule) Install() error {
	// Ensure stackhead user can reload nginx
	permissions := "\n%stackhead ALL= NOPASSWD: /bin/systemctl reload nginx\n"
	if err := xfs.AppendToFile("ssh:///etc/sudoers.d/stackhead", permissions, true); err != nil {
		logger.Debugln(err)
		return fmt.Errorf("unable to add Nginx reload permissions for stackhead user")
	}
	// Validate sudoers file
	if _, _, err := system.RemoteRun("/usr/sbin/visudo -cf /etc/sudoers"); err != nil {
		return fmt.Errorf("unable to validate sudoers file")
	}

	//- set_fact:
	//    nginx_conf_template: "{{ module_role_path | default(role_path) }}/templates/nginx/nginx.conf.j2"

	if err := system.InstallPackage([]system.Package{
		{
			Name:   "nginx",
			Vendor: system.PackageVendorApt,
		},
	}); err != nil {
		return err
	}

	//- name: Setup Nginx
	//  vars:
	//    nginx_ppa_use: true
	//    nginx_vhosts: []
	//    __nginx_user: "stackhead"
	//    root_group: "stackhead"
	//    nginx_extra_conf_options: "{{ module.config.extra_conf_options|default({}) }}"
	//    nginx_extra_conf_http_options: "{{ module.config.extra_conf_http_options|default({}) }}"
	//    server_names_hash_bucket_size: "{{ module.config.server_names_hash_bucket_size|default(64) }}"
	//  include_role:
	//    name: geerlingguy.nginx

	// adjust owner of /var/www directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/var/www"); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-enabled directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/nginx/sites-enabled"); err != nil {
		return err
	}
	// adjust owner of /etc/nginx/sites-available directories
	if _, _, err := system.RemoteRun("chown", "-R", "stackhead:stackhead", "/etc/nginx/sites-available"); err != nil {
		return err
	}

	// Check content after provisioning
	//- name: Check content after provisioning
	//  uri:
	//    url: "http://{{ ansible_default_ipv4.address|default(ansible_all_ipv4_addresses[0]) }}"
	//    return_content: yes
	//  register: uri_result
	//  until: '"Welcome to nginx" in uri_result.content'
	//  retries: 5
	//  delay: 1
	return nil
}
