MAKEFLAGS += --silent
.PHONY: build repl

default: repl

bin/cel-repl: *.go
	go build -o bin/cel-repl

build: bin/cel-repl

repl: build
	./bin/cel-repl
