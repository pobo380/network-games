package model

type Deck struct {
	Cards
}

func (d *Deck) Draw() (Card, bool) {
	c, exist := d.Top()
	d.RemoveTop()
	return c, exist
}
