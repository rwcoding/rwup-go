package internal

type PropertySort []Property

func (ps PropertySort) Len() int {
	return len(ps)
}

func (ps PropertySort) Less(i, j int) bool {
	return ps[i].Order < ps[j].Order
}

func (ps PropertySort) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}
