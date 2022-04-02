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
