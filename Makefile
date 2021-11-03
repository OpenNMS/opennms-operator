.PHONY: build unit-test integration-test validate-build

all: build

unit-test:
	echo "These are unit test"

integration-test:
	echo "These are integration tests"

validate-build:
	circleci config validate .circleci/config.yml

build:
	echo "Hello world!"