class FilterModule(object):
    """
    Jinja2 filters for container related stuff
    """

    def filters(self):
        return {
            'containerPortsStr': self.container_ports_str
        }

    def container_ports_str(self, containerapp__expose, project_name, split_char):
        containerapp__expose.sort(key=lambda x: x['service'])
        output = []

        processed = []
        previous_service = ""
        index = 0
        for nginx_expose in containerapp__expose:
            service_name = nginx_expose['service']
            identifier = service_name + '-' + str(nginx_expose['internal_port'])
            if service_name is None or nginx_expose['external_port'] == 443 or identifier in processed:
                continue
            if previous_service != service_name:
                index = 0
            output.append("${docker_container.stackhead-" + project_name + "-" + service_name + ".ports[" + str(index) + "].external}")
            previous_service = service_name
            index += 1
            processed.append(identifier)

        return split_char.join(output)
