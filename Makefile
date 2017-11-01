GOBIN ?= $(shell go env GOBIN)

help:
	@echo "make run-console          Starts Console http server"
	@echo "make setup-dev            Setup local dev environment"
	@echo "---"
	@echo "make check                Runs all acceptance checks"
	@echo "make check-dependencies   Validates that vendored dependencies are satisfied"
	@echo "---"
	@echo "make test                 Run all tests"
	@echo "make test-integration     Run all integration tests"
	@echo "---"
	@echo "make ui-bundle            Bundle static assets into Go source code (requires esc)"
	@echo "make ui-compile           Compiles Elm code to Javascript (requires elm-make)"
	@echo "make ui-compile-watch     Recompiles the Elm code on file changes (requires elm-live)"

run-config:
	go run ./cmd/configsum/*.go config

run-console:
	cd ui && go run ../cmd/configsum/*.go console -ui.local

setup-dev:
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_dev'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_dev"
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_test'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_test"

dependencies: $(GOBIN)/dep
	$(GOBIN)/dep ensure

check: check-dependencies

check-dependencies: $(GOBIN)/dep
	$(GOBIN)/dep status

test: test-integration

test-integration:
	go test -tags integration ./...

ui-bundle:
	cd ui && make bundle

ui-compile:
	cd ui && make compile

ui-compile-watch:
	cd ui && make compile-watch

.PHONY: setup-dev test test-integration ui-bundle

$(GOBIN)/dep:
	go get -u github.com/golang/dep/cmd/dep