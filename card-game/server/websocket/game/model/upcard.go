package model

type Upcards struct {
	Cards
}

func (up *Upcards) Put(c Card) {
	up.Cards = append(up.Cards, c)
}
