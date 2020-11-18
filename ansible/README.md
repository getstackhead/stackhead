# StackHead Ansible collection

StackHead is aiming to become a container-based Open Source Web Server Management software.

Please visit the repository and documentation for more information on StackHead and how to use it.

* Website: [https://stackhead.io](https://stackhead.io)
* Documentation: [https://docs.stackhead.io](https://docs.stackhead.io)

## Tests

### Unit tests

Due to the collision of the "ansible" python package and our local "ansible" folder,
we have to back one directory in order to run tests.

```
(cd .. && python -m unittest stackhead/ansible/test/units/plugins/action/module_load_config.py)
```
