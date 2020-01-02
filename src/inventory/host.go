package inventory

type Host struct {
	Name string
}

func NewHost() *Host {
	h := &Host{}
	return h
}
