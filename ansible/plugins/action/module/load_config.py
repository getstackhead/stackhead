from ansible.plugins.action import ActionBase
from ansible.config.manager import ensure_type
from ansible.module_utils.six import string_types
from ansible.errors import AnsibleActionFail
from os import path
import re
import yaml


def module_vars_name(module_name):
    if not isinstance(module_name, str):
        return module_name

    chars = re.compile(r'[\\._-]')
    module_name = chars.sub("", module_name)

    return 'module_vars_' + module_name


def tf_populate_module_config(included_module_config, module_rolepath):
    if 'terraform' not in included_module_config:
        return included_module_config
    if 'provider' not in included_module_config['terraform']:
        return included_module_config
    if 'init' not in included_module_config['terraform']['provider']:
        return included_module_config

    init_path = included_module_config['terraform']['provider']['init']

    init_path = re.compile(r'^{{\s*role_path.*\s*}}(.*)$').sub(module_rolepath + "\\1", init_path)
    init_path = re.compile(r'^{{\s*module_role_path.*\s*}}(.*)$').sub(module_rolepath + "\\1", init_path)

    included_module_config['terraform']['provider']['init'] = init_path
    return included_module_config


def _get_role_path(role_name, collection_paths, role_paths):
    parts = role_name.split('.', 2)
    if len(parts) == 3:  # collection
        search_paths = list(map(
            lambda x: x + "/ansible_collections/" + parts[0] + "/" + parts[1] + "/roles",
            collection_paths
        ))
        role_name = parts[2]
    else:  # role
        search_paths = role_paths

    for spath in search_paths:
        full_path = spath + "/" + role_name
        if path.isdir(full_path):
            return full_path

    return None


def load_config(module_name, collection_paths, role_paths):
    module_rolepath = _get_role_path(module_name, collection_paths, role_paths)
    if not module_rolepath:
        raise ValueError("Unable to get role path for module " + module_name)
    module_config_path = module_rolepath + '/stackhead-module.yml'

    config_file = open(module_config_path)
    config_content = yaml.safe_load(config_file)

    populated_config = tf_populate_module_config(config_content, module_rolepath)

    populated_config['role_path'] = module_rolepath

    include_varname = module_vars_name(module_name)
    return {
        'config': populated_config,
        'role_path': module_rolepath,
        'config_varname': include_varname,
        'ansible_facts': {
            include_varname: populated_config
        }
    }


class ActionModule(ActionBase):
    def run(self, tmp=None, task_vars=None):
        super(ActionModule, self).run(tmp, task_vars)
        self._validateOptions()

        return load_config(
            module_name=self._task.args.get('name', None),
            collection_paths=self._lookup_collection_paths(),
            role_paths=self._lookup_role_paths()
        )

    def _lookup_collection_paths(self):
        return self._templar.template("{{ lookup('config', 'COLLECTIONS_PATHS')}}")

    def _lookup_role_paths(self):
        return self._templar.template("{{ lookup('config', 'DEFAULT_ROLES_PATH')}}")

    def _validateOptions(self):
        # Options type validation
        # stings
        for s_type in ('name'):
            if s_type in self._task.args:
                value = ensure_type(self._task.args[s_type], 'string')
                if value is not None and not isinstance(value, string_types):
                    raise AnsibleActionFail("%s is expected to be a string, but got %s instead" % (s_type, type(value)))
                self._task.args[s_type] = value
