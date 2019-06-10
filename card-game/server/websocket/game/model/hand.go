package model

type Hand struct {
	Cards
}

func (h *Hand) Discard(c Card) bool {
	return h.Cards.Remove(c)
}
