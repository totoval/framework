package bigint

type BigIntSlice []BigInt

func (b BigIntSlice) Len() int {
	return len(b)
}

func (b BigIntSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b BigIntSlice) Less(i, j int) bool {
	if b[i].Cmp(b[j]) < 0 {
		return true
	}
	return false
}
