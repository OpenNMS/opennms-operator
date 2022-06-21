#!/usr/bin/env bash

helm repo add horizon-stream https://opennms.github.io/horizon-stream/charts/packaged
helm upgrade -i operator-dependencies horizon-stream/opennms-operator-dependencies