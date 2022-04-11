package alloc

type BufPool [][]byte

func newBufs(size int) BufPool {
	return make([][]byte, size)
}

func (b BufPool) Reset() {
	for i := 0; i < len(b); i++ {
		b[i] = nil
	}
}

type Buffer struct {
	buf    []byte
	offset int
}

func NewBuffer(size int) *Buffer {
	buf := make([]byte, size)
	for i := 0; i < size; i += 4096 {
		buf[i] = 0
	}

	return &Buffer{
		buf: buf,
	}
}

func (b *Buffer) WriteString(s string) {
	off := b.offset + len(s)
	copy(b.buf[b.offset:], s)
	b.offset = off
}

func (b *Buffer) WriteStrings(ss ...string) {
	for i, s := range ss {
		b.WriteString(s)
		if i != len(ss)-1 {
			b.WriteString(",")
		}
	}
	b.WriteString("\n")
}

func (b *Buffer) Bytes() []byte {
	return b.buf[:b.offset]
}
