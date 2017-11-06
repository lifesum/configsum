GOBIN   ?= $(shell go env GOBIN)
NODEBIN ?= $(shell npm bin)

export PATH := $(NODEBIN):$(PATH)

help:
	@echo "make bundle            Bundle static assets into Go source code (requires esc)"
	@echo "make compile           Compiles Elm code to Javascript (requires elm-make)"
	@echo "make compile-watch     Recompiles the Elm code on file changes (requires elm-live)"
	@echo "---"
	@echo "make accept            Run all checks and tests"

bundle: $(GOBIN)/esc
	$(GOBIN)/esc -pkg ui -ignore '.git*|elm-*|node_modules|src|tests|Makefile|package-lock.json' -o ../pkg/ui/static.go -prefix "ui/" -private .

compile: $(NODEBIN)/elm-make
	$(NODEBIN)/elm-make --yes src/Main.elm --output=scripts/console.js

compile-watch: $(NODEBIN)/elm-make $(NODEBIN)/elm-live
	$(NODEBIN)/elm-live --yes src/Main.elm --output=scripts/console.js

# Checks and tests
accept: check

check: check-fmt

check-fmt: $(NODEBIN)/elm-format
	$(NODEBIN)/elm-format --yes src/

.PHONY: accept compile compile-watch static check check-fmt

# Dependencies

$(GOBIN)/esc:
	go get -u github.com/mjibson/esc

$(NODEBIN)/elm-make:
	npm install elm

$(NODEBIN)/elm-live:
	npm install elm-live