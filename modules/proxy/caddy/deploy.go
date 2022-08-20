package proxy_caddy

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3"
	xfs "github.com/saitho/golang-extended-fs"

	"github.com/getstackhead/stackhead/modules/proxy"
	"github.com/getstackhead/stackhead/project"
	"github.com/getstackhead/stackhead/system"
)

type Data struct {
	Project             *project.Project
	AllPortsTfString    string
	DependentContainers string
}

func (CaddyProxyModule) Deploy() error {

	// - name: Include OS-specific variables.
	//  include_vars: "{{ ansible_os_family }}.yml"
	//  ignore_errors: yes
	//
	//- name: Check if authentications are defined
	//  set_fact:
	//    auths_basic: "{{ auths_basic|default([]) + item.security.authentication }}"
	//  when: item.security is defined and item.security.authentication is defined
	//  with_items: "{{ app_config.domains }}"
	//
	//- name: Generate Caddy Terraform file
	//  include_tasks: "{{ module_role_path|default(role_path) }}/tasks/steps/generate-serverconfig-tf.yml"

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

	fmt.Println("Deploy step")
	var buf bytes.Buffer
	fileContents, err := xfs.ReadFile("pkging:///templates/modules/proxy/caddy/caddy.tf.tmpl")
	if err != nil {
		return err
	}
	tmpl, err := template.New("caddy").Funcs(proxy.FuncMap).Funcs(sprig.TxtFuncMap()).Funcs(gomplate.CreateFuncs(nil, nil)).Parse(fileContents)
	if err != nil {
		return err
	}

	data := Data{
		Project:             system.Context.Project,
		AllPortsTfString:    strings.Join(AllPortsTfStrings, ","),
		DependentContainers: strings.Join(DependentContainers, ","),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	err = xfs.WriteFile("ssh://"+system.Context.Project.GetTerraformDirectoryPath()+"/caddy.tf", buf.String())
	return err
}
