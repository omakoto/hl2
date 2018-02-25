package util

type BytesRingBuffer struct {
	capacity int
	start    int
	current  int
	length   int
	buffer   [][]byte
}

func NewStringRingBuffer(capacity int) *BytesRingBuffer {
	if capacity < 0 {
		panic("Negative capacity.")
	}
	return &BytesRingBuffer{capacity: capacity, buffer: make([][]byte, capacity)}
}

func (b *BytesRingBuffer) advance(index *int) {
	if *index < b.capacity-1 {
		*index = *index + 1
	} else {
		*index = 0
	}
}

func (b *BytesRingBuffer) Add(bytes []byte) {
	if b.capacity > 0 {
		b.buffer[b.current] = bytes
		b.advance(&b.current)
		if b.length < b.capacity {
			b.length++
		} else {
			b.advance(&b.start)
		}
	}
}

func (b *BytesRingBuffer) Length() int {
	return b.length
}

func (b *BytesRingBuffer) Clear() {
	b.start = 0
	b.length = 0
	b.current = 0
}

func (b *BytesRingBuffer) For(maxLoop int, f func(bytes []byte)) {
	pos := b.start

	loop := IntMin(maxLoop, b.length)
	for i := 0; i < loop; i++ {
		f(b.buffer[pos])
		b.advance(&pos)
	}
}
