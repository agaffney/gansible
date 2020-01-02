package inventory

import ()

type Inventory struct {
	Hosts  []Host
	Groups []Group
}

func NewInventory() *Inventory {
	i := &Inventory{}
	return i
}
