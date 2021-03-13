package main

import "github.com/google/cel-go/cel"

type historyEntry struct {
	issues *cel.Issues
	ast    *cel.Ast
	raw    string
}

type historyRingBuffer struct {
	size    int
	current int
	buffer  []*historyEntry
}

func newHistoryRingBuffer(size int) *historyRingBuffer {
	buffer := make([]*historyEntry, size)

	return &historyRingBuffer{
		size:    size,
		current: 0,
		buffer:  buffer,
	}
}

func (b *historyRingBuffer) insert(entry *historyEntry) {
	b.buffer[b.position()] = entry
	b.current++
}

func (b *historyRingBuffer) get(index int) *historyEntry {
	if index >= b.size || index < 0 {
		return nil
	}

	return b.buffer[index]
}

func (b *historyRingBuffer) position() int {
	return b.current%b.size
}