from __future__ import absolute_import

import unittest
from unittest.mock import patch

# pylint: disable = relative-beyond-top-level
from .....plugins.action.module.load_config import (
    get_role_path,
    tf_populate_module_config,
)


class TestCopyResultExclude(unittest.TestCase):
    def test_replacing_role_path_in_terraform_provider_init_path(self):
        # No terraform variable
        result = tf_populate_module_config({"foo": "bar"}, "/some/path")
        self.assertEqual({"foo": "bar"}, result)

        # No terraform provider variable
        result = tf_populate_module_config(
            {"foo": "bar", "terraform": {}}, "/some/path"
        )
        self.assertEqual({"foo": "bar", "terraform": {}}, result)

        # No terraform provider init variable
        result = tf_populate_module_config(
            {"foo": "bar", "terraform": {"provider": {}}}, "/some/path"
        )
        self.assertEqual({"foo": "bar", "terraform": {"provider": {}}}, result)

        # No role_path variable in init
        result = tf_populate_module_config(
            {"foo": "bar", "terraform": {"provider": {"init": "test"}}}, "/some/path"
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "test"}}}, result
        )

        # role_path variable in init
        result = tf_populate_module_config(
            {"foo": "bar", "terraform": {"provider": {"init": "{{ role_path }}/test"}}},
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

        # role_path variable in init
        result = tf_populate_module_config(
            {"foo": "bar", "terraform": {"provider": {"init": "{{role_path}}/test"}}},
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

        # module_role_path variable in init
        result = tf_populate_module_config(
            {
                "foo": "bar",
                "terraform": {"provider": {"init": "{{ module_role_path }}/test"}},
            },
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

        # module_role_path variable in init
        result = tf_populate_module_config(
            {
                "foo": "bar",
                "terraform": {"provider": {"init": "{{module_role_path}}/test"}},
            },
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

        # module_role_path|default(role_path) variable in init
        result = tf_populate_module_config(
            {
                "foo": "bar",
                "terraform": {
                    "provider": {
                        "init": "{{ module_role_path | default(role_path) }}/test"
                    }
                },
            },
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

        # module_role_path|default(role_path) variable in init
        result = tf_populate_module_config(
            {
                "foo": "bar",
                "terraform": {
                    "provider": {"init": "{{module_role_path|default(role_path)}}/test"}
                },
            },
            "/some/path",
        )
        self.assertEqual(
            {"foo": "bar", "terraform": {"provider": {"init": "/some/path/test"}}},
            result,
        )

    def test_get_role_path_of_regular_role(self):
        with patch("os.path.isdir") as mocked_isdir:
            mocked_isdir.return_value = True
            result = get_role_path(
                "foo.bar", ["/some/collection/path"], ["/some/role/path"]
            )
            self.assertEqual("/some/role/path/foo.bar", result)

    def test_get_role_path_of_collection_role(self):
        with patch("os.path.isdir") as mocked_isdir:
            mocked_isdir.return_value = True
            result = get_role_path(
                "foo.bar.myrole", ["/some/collection/path"], ["/some/role/path"]
            )
            self.assertEqual(
                "/some/collection/path/ansible_collections/foo/bar/roles/myrole", result
            )
