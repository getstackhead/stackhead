class FilterModule(object):
    """
    Jinja2 filters for container related stuff
    """

    def filters(self):
        return {
            'containerPorts': self.container_ports
        }

    def container_ports(self, containerapp__expose, project_name):
        containerapp__expose.sort(key=lambda x: x['service'])
        output = []

        processed = []
        previous_service = ""
        index = 0
        for nginx_expose in containerapp__expose:
            service_name = nginx_expose['service']
            internal_port = nginx_expose['internal_port']
            identifier = service_name + '-' + str(internal_port)
            if service_name is None or nginx_expose['external_port'] == 443 or identifier in processed:
                continue
            if previous_service != service_name:
                index = 0
            output.append({
                'index': len(output),
                'service': service_name,
                'internal_port': internal_port,
                'tfstring': "${docker_container.stackhead-" + project_name + "-" +
                            service_name + ".ports[" + str(index) + "].external}"
            })
            previous_service = service_name
            index += 1
            processed.append(identifier)

        return output
