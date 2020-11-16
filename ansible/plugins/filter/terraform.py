import re


class FilterModule(object):
    """
    Jinja2 filters for replacing texts by Terraform variables
    """

    def filters(self):
        return {
            'TFreplace': self.tf_replace,
            'TFescapeDoubleQuotes': self.tf_escape_double_quotes,
            'TFpopulateModuleConfig': self.tf_populate_module_config
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

        # already escaped double quotes will get another backslash \" => \\", \\" => \\\"
        double_quotes = re.compile(r'(\\+\")')
        text = double_quotes.sub("\\\\\\1", text)

        # all double quotes have to be escaped with one backslashes again
        text = text.replace('"', '\\"')

        return text

    def tf_populate_module_config(self, included_module_config, module_rolepath):


        init_path = included_module_config.terraform.provider.init

        double_quotes = re.compile(r'^{{ role_path }}(.*)$')
        init_path = double_quotes.sub(module_rolepath + "\1", init_path)

        included_module_config.terraform.provider.init = init_path
        return included_module_config
