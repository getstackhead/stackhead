import re

regex = r"(.*)\.(.*\..*)"


def domain(input_str):
    matches = re.search(regex, input_str)
    if not matches:
        return None
    return matches.group(2)


def subdomain(input_str):
    matches = re.search(regex, input_str)
    if not matches:
        return None
    return matches.group(1)


class FilterModule(object):
    filter_map = {
        'domain': domain,
        'subdomain': subdomain,
    }

    def filters(self):
        return self.filter_map
