package main

import (
	"log"
)

func main() {
	repl, err := NewRepl()

	if err != nil {
		log.Fatal("Failed to instantiate repl", err)
	}
	repl.init()
	repl.loop()
}
