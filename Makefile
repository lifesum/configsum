help:
	@echo "make setup-dev         Setup local dev environment"
	@echo "--"
	@echo "make test              Run all tests"
	@echo "make test-integration  Run all integration tests"

run-console:
	cd ui && go run ../cmd/configsum/*.go console -static.local

setup-dev:
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_dev'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_dev"
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_test'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_test"

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
