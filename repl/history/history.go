package history

func NewEntryRingBuffer(size int) *EntryRingBuffer {
	buffer := make([]*Entry, size)

	return &EntryRingBuffer{
		size:    size,
		current: 0,
		buffer:  buffer,
	}
}

func (b *EntryRingBuffer) Insert(entry *Entry) {
	b.buffer[b.Position()] = entry
	b.current++
}

func (b *EntryRingBuffer) Get(index int) *Entry {
	if index >= b.size || index < 0 {
		return nil
	}

	return b.buffer[index]
}

func (b *EntryRingBuffer) Position() int {
	return b.current%b.size
}
