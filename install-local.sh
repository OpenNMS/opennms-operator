#!/usr/bin/env bash
#use this script to install a basic version of the opennms operator locally

make local-docker

helm upgrade -i operator-local ./charts/opennms-operator -f local.yaml