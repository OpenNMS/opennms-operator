.PHONY: dependencies build alpine-build unit-test integration-test validate-build
.PHONY: helm-dep helm-package

all: build

unit-test:
	echo "These are unit test"

integration-test:
	echo "These are integration tests"

validate-build:
	circleci config validate .circleci/config.yml

helm-dep: charts/opennms-operator/charts
	helm dep update charts/opennms-operator

helm-package: helm-dep
	helm package charts/opennms-operator
	helm repo index --url https://opennms.github.io/opennms-operator/ --merge index.yaml .

dependencies:
	go mod download

build:
	go build -a -o operator cmd/opennms-operator/main.go

alpine-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o operator cmd/opennms-operator/main.go