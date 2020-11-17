import re

class FilterModule(object):
    """
    Jinja2 filters for modules
    """

    def filters(self):
        return {
            'moduleVarsName': self.module_vars_name
        }

    def module_vars_name(self, module_name):
        if not isinstance(module_name, str):
            return module_name

        chars = re.compile(r'[\\._-]')
        module_name = chars.sub("", module_name)

        return 'module_vars_' + module_name
