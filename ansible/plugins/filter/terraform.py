import re


class FilterModule(object):
    """
    Jinja2 filters for replacing texts by Terraform variables
    """

    def filters(self):
        return {
            "TFreplace": self.tf_replace,
            "TFescapeDoubleQuotes": self.tf_escape_double_quotes,
        }

    def tf_replace(
        self, text, project_name, container_resource_name="docker_container"
    ):
        if not isinstance(text, str):
            return text
        # Replace Docker service name variables
        docker_service = re.compile(r"\$DOCKER_SERVICE_NAME\[\'(.*)\'\]")
        text = docker_service.sub(
            "${"
            + container_resource_name
            + ".stackhead-"
            + project_name
            + "-\\1.name}",
            text,
        )

        return text

    def tf_escape_double_quotes(self, text):
        if not isinstance(text, str):
            return text

        # duplicate existing escapes \ => \\, \" => \\", \\" => \\\\"
        # and prepend one additional escape \\ => \\\, \\\\" => \\\\\", " => \"
        text = re.compile(r"((\\*)\"|(\\+))").sub("\\\\\\3\\2\\1", text)

        return text
