.PHONY: build unit-test integration-test

all: build

unit-test:
	echo "These are unit test"

integration-test:
	echo "These are integration tests"

build: test
	echo "Hello world!"