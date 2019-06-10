package action

import (
	"reflect"
	"testing"

	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

func TestPass_Do(t *testing.T) {
	type fields struct {
		PlayerId model.PlayerId
	}
	type args struct {
		st *state.State
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantEvents   event.Events
		wantRetState *state.State
	}{
		{
			name: "PlayerId is not found",
			fields: fields{
				PlayerId: "3",
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
				},
			},
			wantEvents: event.Events{
				&event.NotAcceptable{
					ActionType: TypePass,
					Reason:     "PlayerId is not found.",
				},
			},
			wantRetState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
			},
		},
		{
			name: "Not turn owner",
			fields: fields{
				PlayerId: "2",
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 0,
					},
				},
			},
			wantEvents: event.Events{
				&event.NotAcceptable{
					ActionType: TypePass,
					Reason:     "PlayerId is not turn owner.",
				},
			},
			wantRetState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 0,
				},
			},
		},
		{
			name: "デッキがなくなって引き分け",
			fields: fields{
				PlayerId: "1",
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 0,
					},
					Deck:     model.Deck{},
					Discards: model.Discards{},
				},
			},
			wantEvents: event.Events{
				&event.Pass{
					PlayerId: "1",
				},
				&event.FinishGame{
					FinishType: event.FinishTypeDraw,
				},
			},
			wantRetState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 0,
				},
				Deck:     model.Deck{},
				Discards: model.Discards{},
			},
		},
		{
			name: "デッキにカードがあるパターン",
			fields: fields{
				PlayerId: "1",
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1", Hand: model.Hand{Cards: model.Cards{{Suit: model.SuitClub, Number: 1}}}},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 0,
					},
					Deck: model.Deck{
						Cards: model.Cards{
							{Suit: model.SuitSpade, Number: 1},
							{Suit: model.SuitSpade, Number: 2},
						},
					},
					Discards: model.Discards{},
				},
			},
			wantEvents: event.Events{
				&event.Pass{
					PlayerId: "1",
				},
				&event.DrawCard{
					PlayerId: "1",
					Cards:    model.Cards{{Suit: model.SuitSpade, Number: 2}},
					CardsNum: 1,
				},
				&event.TurnChange{
					PlayerId: "2",
				},
			},
			wantRetState: &state.State{
				Players: []model.Player{
					{
						Id: "1",
						Hand: model.Hand{
							Cards: model.Cards{
								{Suit: model.SuitClub, Number: 1},
								{Suit: model.SuitSpade, Number: 2},
							},
						},
					},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 1,
				},
				Deck: model.Deck{
					Cards: model.Cards{
						{Suit: model.SuitSpade, Number: 1},
					},
				},
				Discards: model.Discards{},
			},
		},
		{
			name: "デッキにカードがないからシャッフルして引くパターン",
			fields: fields{
				PlayerId: "1",
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{
							Id: "1",
							Hand: model.Hand{
								Cards: model.Cards{
									{Suit: model.SuitClub, Number: 1},
								},
							},
						},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 0,
					},
					Deck: model.Deck{
						Cards: model.Cards{},
					},
					Discards: model.Discards{
						Cards: model.Cards{
							{Suit: model.SuitSpade, Number: 1},
						},
					},
				},
			},
			wantEvents: event.Events{
				&event.Pass{
					PlayerId: "1",
				},
				&event.DrawCard{
					PlayerId: "1",
					Cards:    model.Cards{{Suit: model.SuitSpade, Number: 1}},
					CardsNum: 1,
				},
				&event.TurnChange{
					PlayerId: "2",
				},
			},
			wantRetState: &state.State{
				Players: []model.Player{
					{
						Id: "1",
						Hand: model.Hand{
							Cards: model.Cards{
								{Suit: model.SuitClub, Number: 1},
								{Suit: model.SuitSpade, Number: 1},
							},
						},
					},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 1,
				},
				Deck: model.Deck{
					Cards: model.Cards{},
				},
				Discards: model.Discards{
					Cards: model.Cards{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pass{
				PlayerId: tt.fields.PlayerId,
			}
			gotEvents, gotRetState := p.Do(tt.args.st)
			if !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("Pass.Do() gotEvents = %+v, want %+v", gotEvents, tt.wantEvents)
			}
			if !reflect.DeepEqual(gotRetState, tt.wantRetState) {
				t.Errorf("Pass.Do() gotRetState = %+v, want %+v", gotRetState, tt.wantRetState)
			}
		})
	}
}
