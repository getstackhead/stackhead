from ..action.module.load_config import module_vars_name


class FilterModule(object):
    """
    Jinja2 filters for module related stuff
    """

    def filters(self):
        return {"module_vars_name": module_vars_name}
