package history

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistoryRingBuffer_Overflow(t *testing.T) {
	b := NewEntryRingBuffer(3)

	b.Insert(&Entry{Raw: "0"})
	b.Insert(&Entry{Raw: "1"})
	b.Insert(&Entry{Raw: "2"})
	b.Insert(&Entry{Raw: "3"})

	assert.Equal(t, len(b.buffer), 3)
	assert.Equal(t, b.Get(0).Raw, "3")
	assert.Equal(t, b.Get(1).Raw, "1")
	assert.Equal(t, b.Get(2).Raw, "2")
}

func TestHistoryRingBuffer_GetOutOfBounds(t *testing.T) {
	b := NewEntryRingBuffer(1)

	b.Insert(&Entry{Raw: "0"})

	assert.Nil(t, b.Get(1))
}

func TestHistoryRingBuffer_GetNegativeIndex(t *testing.T) {
	b := NewEntryRingBuffer(1)

	b.Insert(&Entry{Raw: "0"})

	assert.Nil(t, b.Get(-1))
}
