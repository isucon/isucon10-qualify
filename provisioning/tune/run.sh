#!/bin/bash
ansible-playbook -i hosts.yaml db.yaml
ansible-playbook -i hosts.yaml app.yaml
ansible-playbook -i hosts.yaml bench.yaml
