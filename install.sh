#!/usr/bin/env bash
#use this script to install a basic version of the opennms operator

# cert-manager
helm repo add jetstack https://charts.jetstack.io

helm upgrade -i \
  onms-cert-manager jetstack/cert-manager \
  --namespace onms-cert-manager \
  --create-namespace \
  --version v1.7.0 \
  --set installCRDs=true

# operator
helm repo add opennms-operator https://opennms.github.io/opennms-operator/charts/packaged

helm upgrade -i opennms-operator opennms-operator/opennms-operator -f values.yaml

# k8s replicator
helm repo add mittwald https://helm.mittwald.de

helm upgrade -i onms-kubernetes-replicator --namespace opennms mittwald/kubernetes-replicator