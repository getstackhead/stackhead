class FilterModule(object):
    """
    Jinja2 filters for Ruby related stuff
    """

    def filters(self):
        return {
            'to_ruby_hash': self.to_ruby_hash
        }

    def to_ruby_hash(self, dict):
        return str(dict).replace(': ', ' => ')

