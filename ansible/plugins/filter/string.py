def endswith(string, suffix):
    return string.endswith(suffix)

class FilterModule(object):
    """
    Jinja2 filters for string related stuff
    """

    def filters(self):
        return {"endswith": endswith}
