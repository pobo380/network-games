package model

import (
	"math/rand"
)

const (
	SuitSpade Suit = iota
	SuitClub
	SuitDiamond
	SuitHeart
)

const (
	JokerNumber       Number = 14
	InvalidCardNumber Number = -1
)

var (
	Suits   = []Suit{SuitSpade, SuitClub, SuitDiamond, SuitHeart}
	Numbers = []Number{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, JokerNumber}

	InvalidCard = Card{Number: InvalidCardNumber}
)

type Suit int
type Number int

type Card struct {
	Number Number
	Suit   Suit
}

func (c *Card) IsJoker() bool {
	return c.Number == JokerNumber
}

func (c *Card) IsInvalid() bool {
	return c.Number == InvalidCardNumber
}

func (c *Card) SameSuit(o Card) bool {
	return c.Suit == o.Suit
}

func (c *Card) SameNumber(o Card) bool {
	return c.Number == o.Number
}

type Cards []Card

func (cs Cards) Shuffle(r *rand.Rand) {
	r.Shuffle(len(cs), func(i, j int) {
		cs[i], cs[j] = cs[j], cs[i]
	})
}

func (cs Cards) Top() (Card, bool) {
	if l := len(cs); l > 0 {
		return cs[l-1], true
	}

	return InvalidCard, false
}

func (cs *Cards) Contain(target Card) bool {
	for _, c := range *cs {
		if c == target {
			return true
		}
	}

	return false
}

func (cs Cards) Empty() bool {
	return len(cs) == 0
}

func (cs *Cards) Add(target Card) {
	*cs = append(*cs, target)
}

func (cs *Cards) Remove(target Card) bool {
	for i, c := range *cs {
		if c == target {
			*cs = append((*cs)[:i], (*cs)[i+1:]...)
			return true
		}
	}

	return false
}

func (cs *Cards) RemoveTop() {
	l := len(*cs)
	if l > 0 {
		*cs = (*cs)[:l-1]
	}
}

func (cs *Cards) Clear() {
	*cs = Cards{}
}
