#!/usr/bin/env bash
#use this script to install a basic version of the opennms operator

# operator
helm repo add opennms-operator https://opennms.github.io/opennms-operator/charts/packaged

helm upgrade -i opennms-operator opennms-operator/opennms-operator -f values.yaml