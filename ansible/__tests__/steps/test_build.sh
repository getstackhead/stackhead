#!/bin/bash
ansible-galaxy collection build -f ansible
ansible-galaxy collection install getstackhead-stackhead-*
