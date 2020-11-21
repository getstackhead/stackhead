import re


class FilterModule(object):
    """
    Jinja2 filters for replacing texts by Terraform variables
    """

    def filters(self):
        return {
            'TFreplace': self.tf_replace,
            'TFescapeDoubleQuotes': self.tf_escape_double_quotes
        }

    def tf_replace(self, text, project_name):
        if not isinstance(text, str):
            return text
        # Replace Docker service name variables
        docker_service = re.compile(r'\$DOCKER_SERVICE_NAME\[\'(.*)\'\]')
        text = docker_service.sub("${docker_container.stackhead-" + project_name + "-\\1.name}", text)

        return text

    def tf_escape_double_quotes(self, text):
        if not isinstance(text, str):
            return text

        # duplicate existing escapes \" => \\", \\" => \\\\"
        double_quotes = re.compile(r'((\\+)\")')
        text = double_quotes.sub("\\2\\1", text)

        # all double quotes have to be escaped with one backslashes again
        text = text.replace('"', '\\"')

        return text
