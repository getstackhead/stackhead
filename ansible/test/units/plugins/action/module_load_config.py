from __future__ import absolute_import

import os
import unittest
from unittest.mock import MagicMock, Mock, patch
from ansible.playbook.task import Task
from ansible.plugins.loader import connection_loader
from .....plugins.action.module.load_config import ActionModule


class TestCopyResultExclude(unittest.TestCase):

    def setUp(self):
        self.play_context = Mock()
        self.play_context.shell = 'sh'
        self.connection = connection_loader.get('local', self.play_context, os.devnull)

    def tearDown(self):
        pass

    def test_replacing_role_path_in_terraform_provider_init_path(self):
        task = MagicMock(Task)
        task.async_val = False
        # task.args = {'_raw_params': 'Args1'}

        self.play_context.check_mode = False

        subject = ActionModule(task, self.connection, self.play_context, loader=None, templar=None,
                               shared_loader_obj=None)

        # No terraform variable
        result = subject.tf_populate_module_config({'foo': 'bar'}, '/some/path')
        self.assertEqual({'foo': 'bar'}, result)

        # No terraform provider variable
        result = subject.tf_populate_module_config({'foo': 'bar', 'terraform': {}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {}}, result)

        # No terraform provider init variable
        result = subject.tf_populate_module_config({'foo': 'bar', 'terraform': {'provider': {}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {}}}, result)

        # No role_path variable in init
        result = subject.tf_populate_module_config({'foo': 'bar', 'terraform': {'provider': {'init': 'test'}}},
                                                   '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': 'test'}}}, result)

        # role_path variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{ role_path }}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

        # role_path variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{role_path}}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

        # module_role_path variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{ module_role_path }}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

        # module_role_path variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{module_role_path}}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

        # module_role_path|default(role_path) variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{ module_role_path | default(role_path) }}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

        # module_role_path|default(role_path) variable in init
        result = subject.tf_populate_module_config(
            {'foo': 'bar', 'terraform': {'provider': {'init': '{{module_role_path|default(role_path)}}/test'}}}, '/some/path')
        self.assertEqual({'foo': 'bar', 'terraform': {'provider': {'init': '/some/path/test'}}}, result)

    def test_get_role_path_of_regular_role(self):
        task = MagicMock(Task)
        task.async_val = False

        self.play_context.check_mode = False

        self.mock_am = ActionModule(task, self.connection, self.play_context, loader=None, templar=None,
                                    shared_loader_obj=None)
        self.mock_am._lookup_role_paths = lambda: ['/some/role/path']

        with patch('os.path.isdir') as mocked_isdir:
            mocked_isdir.return_value = True
            result = self.mock_am._get_role_path('foo.bar')
            self.assertEqual('/some/role/path/foo.bar', result)

    def test_get_role_path_of_collection_role(self):
        task = MagicMock(Task)
        task.async_val = False

        self.play_context.check_mode = False

        self.mock_am = ActionModule(task, self.connection, self.play_context, loader=None, templar=None,
                                    shared_loader_obj=None)
        self.mock_am._lookup_collection_paths = lambda: ['/some/collection/path']

        with patch('os.path.isdir') as mocked_isdir:
            mocked_isdir.return_value = True
            result = self.mock_am._get_role_path('foo.bar.myrole')
            self.assertEqual('/some/collection/path/ansible_collections/foo/bar/roles/myrole', result)

