package inventory

type Group struct {
	Name string
}

func NewGroup() *Group {
	g := &Group{}
	return g
}
