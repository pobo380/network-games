package state

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"math/rand"
	"time"
)

type State struct {
	// Config
	rand   *rand.Rand
	Config model.Config

	// GameMaster
	PlayOrder model.PlayOrder

	// Player
	Players []model.Player

	// Table
	Deck     model.Deck
	Discards model.Discards
	Upcards  model.Upcards
}

//
// State
//

func NewState(config model.Config, players []model.Player) *State {
	return &State{
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		Config:  config,
		Players: players,
	}
}

func (st *State) InitGame() {
	// init deck/discards/upcards
	st.InitDeck()
	st.Discards.Cards = model.Cards{}
	st.Upcards.Cards = model.Cards{}

	// deal to player hands
	st.Deck.Shuffle(st.rand)
	st.Deal()

	// determine PlayOrder
	st.PlayOrder = model.PlayOrder{
		CurrentIdx: 0,
	}
	for _, player := range st.Players {
		st.PlayOrder.Order = append(st.PlayOrder.Order, player.Id)
	}

	st.PlayOrder.Shuffle(st.rand)

	// set initial upcards
	st.Upcards.Cards = model.Cards{}
	c, _ := st.Deck.Draw()
	st.Upcards.Cards.Add(c)
}

func (st *State) InitDeck() {
	for _, suit := range model.Suits {
		for _, number := range model.Numbers {
			st.Deck.Cards = append(st.Deck.Cards, model.Card{Suit: suit, Number: number})
		}
	}
}

func (st *State) Deal() {
	num := st.Config.InitialHandNum

	for i := range st.Players {
		h, t := i*num, (i+1)*num
		st.Players[i].Hand.Cards = st.Deck.Cards[h:t]
	}

	st.Deck.Cards = st.Deck.Cards[len(st.Players)*num:]
}

func (st *State) ResetDeck() {
	st.Deck.Cards = st.Discards.Cards
	st.Deck.Shuffle(st.rand)
	st.Discards.Clear()
}

//
// Player
//

func (st *State) FindPlayerById(id model.PlayerId) *model.Player {
	for i := range st.Players {
		pl := &st.Players[i]
		if id == pl.Id {
			return pl
		}
	}

	return nil
}
