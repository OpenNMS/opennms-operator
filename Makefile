.PHONY: dependencies build alpine-build unit-test integration-test validate-build
.PHONY: helm-dep helm-package

all: build

unit-test:
	go test --tags=unit ./...

integration-test:
	echo "These are integration tests"

validate-build:
	circleci config validate .circleci/config.yml

helm-dep:
	helm dep update charts/opennms-operator

helm-package: helm-dep
	helm package charts/opennms-operator -d charts/packaged
	helm repo index --url https://opennms.github.io/opennms-operator/charts/packaged --merge index.yaml .

dependencies:
	go mod download

build:
	go build -a -o operator cmd/opennms-operator/main.go

alpine-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o operator cmd/opennms-operator/main.go