#!/usr/bin/env bash
#use this script to install a basic version of the opennms operator locally

echo __________________Installing OLM___________________
echo
operator-sdk olm install

echo
echo _______________Building Docker Image_______________
echo
make local-docker


echo
echo ___________Installing Helm Dependencies____________
echo
bash bin/install-helm-deps.sh

echo
echo _______________Wait For Dependencies_______________
echo
until kubectl wait -n kafka --for=condition=Ready=true pod -l name=strimzi-cluster-operator --timeout=90s 2> /dev/null
do
    sleep 5
    echo Waiting for dependencies to start....
done
echo Dependencies started.

echo
echo ________________Installing Operator________________
echo
helm upgrade -i operator-local ./charts/opennms-operator -f values.yaml