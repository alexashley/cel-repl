MAKEFLAGS += --silent
.PHONY: build repl

GO_SRC = $(shell find . -type f -name '*.go')

default: repl

bin/cel-repl: $(GO_SRC)
	go build -o bin/cel-repl

build: bin/cel-repl

repl: build
	./bin/cel-repl
