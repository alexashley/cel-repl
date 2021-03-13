package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistoryRingBuffer_Overflow(t *testing.T) {
	b := newHistoryRingBuffer(3)

	b.insert(&historyEntry{raw: "0"})
	b.insert(&historyEntry{raw: "1"})
	b.insert(&historyEntry{raw: "2"})
	b.insert(&historyEntry{raw: "3"})

	assert.Equal(t, len(b.buffer), 3)
	assert.Equal(t, b.get(0).raw, "3")
	assert.Equal(t, b.get(1).raw, "1")
	assert.Equal(t, b.get(2).raw, "2")
}

func TestHistoryRingBuffer_GetOutOfBounds(t *testing.T) {
	b := newHistoryRingBuffer(1)

	b.insert(&historyEntry{raw: "0"})

	assert.Nil(t, b.get(1))
}

func TestHistoryRingBuffer_GetNegativeIndex(t *testing.T) {
	b := newHistoryRingBuffer(1)

	b.insert(&historyEntry{raw: "0"})

	assert.Nil(t, b.get(-1))
}
