help:
	@echo "make setup-dev         Setup local dev environment"

setup-dev:
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_dev'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_dev"
	psql -d template1 -tc "SELECT 1 FROM pg_database WHERE datname = 'configsum_test'" | grep -q 1 || psql -d template1 -c "CREATE DATABASE configsum_test"
