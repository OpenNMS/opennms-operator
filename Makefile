.PHONY: dependencies build unit-test integration-test validate-build

all: build

unit-test:
	echo "These are unit test"

integration-test:
	echo "These are integration tests"

validate-build:
	circleci config validate .circleci/config.yml

dependencies:
	go mod download

build:
	go build -a -o operator cmd/opennms-operator/main.go

alpine-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o operator cmd/opennms-operator/main.go