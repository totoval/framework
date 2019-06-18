package bigfloat

type Slice []*BigFloat

func (s Slice) Len() int {
	return len(s)
}

func (s Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Slice) Less(i, j int) bool {
	if s[i].Cmp(*s[j]) < 0 {
		return true
	}
	return false
}
