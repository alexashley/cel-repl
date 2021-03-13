package main

import (
	_ "embed"
	"github.com/alexashley/cel-repl/repl"
	"log"
)

//go:embed version
var replVersion string

func main() {
	r, err := repl.NewRepl(replVersion)

	if err != nil {
		log.Fatal("Failed to instantiate repl", err)
	}

	r.Init().Loop()
}
