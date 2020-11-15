from ansible.plugins.lookup import LookupBase
from os import path

class LookupModule(LookupBase):

    def run(self, terms, variables=None, **kwargs):
        self.set_options(direct=kwargs)

        role_name = terms[0]

        parts = role_name.split('.', 2)
        if len(parts) == 3:  # collection
            collection_paths = self._templar.template("{{ lookup('config', 'COLLECTIONS_PATHS')}}")
            search_paths = list(map(lambda x: x + "/ansible_collections/" + parts[0] + "/" + parts[1] + "/roles", collection_paths))
            role_name = parts[2]
        else: # role
            search_paths = self._templar.template("{{ lookup('config', 'DEFAULT_ROLES_PATH')}}")

        ret = []
        for spath in search_paths:
            full_path = spath + "/" + role_name
            if path.isdir(full_path):
                ret.append(full_path)

        return ret
