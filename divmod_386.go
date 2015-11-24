package base58

// could use multiply and shift since the compiler doesn't recognize
// this sequence.
func divmod(a uint64) (q, r uint64) {
	q = a / 58
	r = a % 58
	return
}
