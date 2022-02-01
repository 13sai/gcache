package gcache

func cloneBytes(b []byte) []byte {
	n := make([]byte, len(b))
	copy(n, b)
	return n
}
