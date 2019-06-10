package model

type Discards struct {
	Cards
}

func (up *Discards) Put(c Card) {
	up.Cards = append(up.Cards, c)
}
