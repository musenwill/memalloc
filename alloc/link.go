package alloc

type Node struct {
	buf  []byte
	next *Node
}

type Link struct {
	head  *Node
	count int64
}

func (ln *Link) Push(buf []byte) {
	node := &Node{
		buf:  buf,
		next: ln.head,
	}
	ln.head = node
	ln.count++
}

func (ln *Link) Reset() {
	ln.head = nil
	ln.count = 0
}

func (ln *Link) Len() int64 {
	return ln.count
}
