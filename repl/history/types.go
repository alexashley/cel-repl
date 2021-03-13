package history

import "github.com/google/cel-go/cel"

type Entry struct {
	Issues *cel.Issues
	Ast    *cel.Ast
	Raw    string
}

type EntryRingBuffer struct {
	size    int
	current int
	buffer  []*Entry
}
