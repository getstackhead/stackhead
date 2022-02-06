package proxy_nginx

import (
	"github.com/getstackhead/stackhead/system"
)

func (NginxProxyModule) Install() error {
	// - name: Ensure stackhead user can reload nginx
	//  blockinfile:
	//    path: '/etc/sudoers.d/stackhead'
	//    block: |
	//      %stackhead ALL= NOPASSWD: /bin/systemctl reload nginx
	//    mode: 0440
	//    create: yes
	//    state: present
	//    validate: 'visudo -cf %s'

	//- name: Deploy additional h5bp Nginx files
	//  synchronize:
	//    src: "{{ module_role_path | default(role_path) }}/vendor/server-configs-nginx/h5bp"
	//    dest: /etc/nginx
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

	//- name: adjust owner of /var/www directories
	//  file:
	//    path: /var/www
	//    state: directory
	//    owner: "stackhead"
	//    group: "stackhead"
	//    mode: 0755
	//    recurse: true
	//- name: adjust owner of /etc/nginx/sites-enabled directory
	//  file:
	//    path: /etc/nginx/sites-enabled
	//    state: directory
	//    owner: "stackhead"
	//    group: "stackhead"
	//    mode: 0755
	//    recurse: true
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
