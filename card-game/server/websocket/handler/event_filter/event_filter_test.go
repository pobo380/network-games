package event_filter

import (
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
	"reflect"
	"testing"

	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
)

func TestFilter(t *testing.T) {
	type args struct {
		src      event.Events
		playerId string
	}
	tests := []struct {
		name    string
		args    args
		wantRet event.Events
	}{
		{
			name: "フィルターされるやつ",
			args: args{
				src: event.Events{
					&event.DrawCard{
						PlayerId: "2",
						Cards:    model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
						CardsNum: 2,
					},
					&event.GameState{
						State: &state.State{
							Players: []model.Player{
								{
									Id: "1",
									Hand: model.Hand{
										Cards: model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
									},
								},
								{
									Id: "2",
									Hand: model.Hand{
										Cards: model.Cards{{Suit: model.SuitDiamond, Number: 1}, {Suit: model.SuitDiamond, Number: 2}},
									},
								},
							},
						},
					},
				},
				playerId: "1",
			},
			wantRet: event.Events{
				&event.DrawCard{
					PlayerId: "2",
					Cards:    model.Cards{},
					CardsNum: 2,
				},
				&event.GameState{
					State: &state.State{
						Players: []model.Player{
							{
								Id: "1",
								Hand: model.Hand{
									Cards: model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
								},
							},
							{
								Id: "2",
								Hand: model.Hand{
									Cards: model.Cards{model.InvalidCard, model.InvalidCard},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := Filter(tt.args.src, tt.args.playerId); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("Filter() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func Test_filterDrawCard(t *testing.T) {
	type args struct {
		ev       event.Event
		playerId string
	}
	tests := []struct {
		name string
		args args
		want event.Event
	}{
		{
			name: "フィルターされる",
			args: args{
				ev: &event.DrawCard{
					PlayerId: "1",
					Cards:    model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
					CardsNum: 2,
				},
				playerId: "2",
			},
			want: &event.DrawCard{
				PlayerId: "1",
				Cards:    model.Cards{},
				CardsNum: 2,
			},
		},
		{
			name: "フィルターされない",
			args: args{
				ev: &event.DrawCard{
					PlayerId: "1",
					Cards:    model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
					CardsNum: 2,
				},
				playerId: "1",
			},
			want: &event.DrawCard{
				PlayerId: "1",
				Cards:    model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
				CardsNum: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterDrawCard(tt.args.ev, tt.args.playerId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterDrawCard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterGameState(t *testing.T) {
	type args struct {
		ev       event.Event
		playerId string
	}
	tests := []struct {
		name string
		args args
		want event.Event
	}{
		{
			name: "フィルターされる",
			args: args{
				ev: &event.GameState{
					State: &state.State{
						Players: []model.Player{
							{
								Id: "1",
								Hand: model.Hand{
									Cards: model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
								},
							},
							{
								Id: "2",
								Hand: model.Hand{
									Cards: model.Cards{{Suit: model.SuitDiamond, Number: 1}, {Suit: model.SuitDiamond, Number: 2}},
								},
							},
						},
					},
				},
				playerId: "1",
			},
			want: &event.GameState{
				State: &state.State{
					Players: []model.Player{
						{
							Id: "1",
							Hand: model.Hand{
								Cards: model.Cards{{Suit: model.SuitClub, Number: 1}, {Suit: model.SuitClub, Number: 2}},
							},
						},
						{
							Id: "2",
							Hand: model.Hand{
								Cards: model.Cards{model.InvalidCard, model.InvalidCard},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterGameState(tt.args.ev, tt.args.playerId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterGameState() = %v, want %v", got, tt.want)
			}
		})
	}
}
