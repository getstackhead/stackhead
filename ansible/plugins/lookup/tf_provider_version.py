from os import popen
import re

from ansible.plugins.lookup import LookupBase


class LookupModule(LookupBase):
    def run(self, terms, variables=None, **kwargs):
        self.set_options(direct=kwargs)

        provider_name = terms[0]  # fully qualified provider name (e.g. getstackhead/acme)
        stream = popen('(cd /stackhead/terraform/projects && terraform providers)')
        output = stream.read()

        matches = re.finditer(r"(├|└)── provider\[.*/(.*/.*)\] (.+)$", output, re.MULTILINE)

        for matchNum, match in enumerate(matches, start=1):
            if match.group(2) == provider_name:
                return [match.group(3)]
        return [""]
