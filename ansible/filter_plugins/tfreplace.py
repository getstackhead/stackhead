import re


class FilterModule(object):
    """
    Jinja2 filters for replacing texts by Terraform variables
    """

    def filters(self):
        return {
            'TFreplace': self.tf_replace
        }

    def tf_replace(self, text, project_name):
        # Replace Docker service name variables
        docker_service = re.compile(r'\$DOCKER_SERVICE_NAME\[\'(.*)\'\]')
        text = docker_service.sub("${docker_container.stackhead-" + project_name + "-\\1.name}", text)

        return text
