#!/bin/bash
cd ansible
ansible-galaxy collection build -f
ansible-galaxy collection install getstackhead-stackhead-*
