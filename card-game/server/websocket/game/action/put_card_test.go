package action

import (
	"reflect"
	"testing"

	"github.com/pobo380/network-games/card-game/server/websocket/game/event"
	"github.com/pobo380/network-games/card-game/server/websocket/game/model"
	"github.com/pobo380/network-games/card-game/server/websocket/game/state"
)

func TestPutCard_Do(t *testing.T) {
	type fields struct {
		playerId model.PlayerId
		card     model.Card
	}
	type args struct {
		st *state.State
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantEvents event.Events
		wantState  *state.State
	}{
		{
			name: "カードが出せて次の手番に回る",
			fields: fields{
				playerId: "1",
				card:     model.Card{Suit: model.SuitClub, Number: 1},
			},
			args: args{
				st: &state.State{
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"2", "1", "3"},
						CurrentIdx: 1,
					},
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
						{
							Id: "2",
						},
						{
							Id: "3",
						},
					},
					Upcards: model.Upcards{
						Cards: model.Cards{
							{Suit: model.SuitClub, Number: 10},
						},
					},
					Discards: model.Discards{
						Cards: model.Cards{},
					},
				},
			},
			wantEvents: event.Events{
				&event.PutCard{
					PlayerId: "1",
					Card:     model.Card{Suit: model.SuitClub, Number: 1},
				},
				&event.TurnChange{
					PlayerId: "3",
				},
			},
			wantState: &state.State{
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"2", "1", "3"},
					CurrentIdx: 2,
				},
				Players: []model.Player{
					{
						Id: "1",
						Hand: model.Hand{
							Cards: model.Cards{
								{Suit: model.SuitSpade, Number: 2},
							},
						},
					},
					{
						Id: "2",
					},
					{
						Id: "3",
					},
				},
				Upcards: model.Upcards{
					Cards: model.Cards{
						{Suit: model.SuitClub, Number: 10},
						{Suit: model.SuitClub, Number: 1},
					},
				},
				Discards: model.Discards{
					Cards: model.Cards{
						{Suit: model.SuitClub, Number: 10},
					},
				},
			},
		},
		{
			name: "カードが出せてゲーム終了",
			fields: fields{
				playerId: "1",
				card:     model.Card{Suit: model.SuitClub, Number: 1},
			},
			args: args{
				st: &state.State{
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"2", "1", "3"},
						CurrentIdx: 1,
					},
					Players: []model.Player{
						{
							Id: "1",
							Hand: model.Hand{
								Cards: model.Cards{
									{Suit: model.SuitClub, Number: 1},
								},
							},
						},
						{
							Id: "2",
						},
						{
							Id: "3",
						},
					},
					Upcards: model.Upcards{
						Cards: model.Cards{
							{Suit: model.SuitClub, Number: 10},
						},
					},
					Discards: model.Discards{
						Cards: model.Cards{},
					},
				},
			},
			wantEvents: event.Events{
				&event.PutCard{
					PlayerId: "1",
					Card:     model.Card{Suit: model.SuitClub, Number: 1},
				},
				&event.FinishGame{
					WinnerId:   "1",
					FinishType: event.FinishTypeWin,
				},
			},
			wantState: &state.State{
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"2", "1", "3"},
					CurrentIdx: 1,
				},
				Players: []model.Player{
					{
						Id: "1",
						Hand: model.Hand{
							Cards: model.Cards{},
						},
					},
					{
						Id: "2",
					},
					{
						Id: "3",
					},
				},
				Upcards: model.Upcards{
					Cards: model.Cards{
						{Suit: model.SuitClub, Number: 10},
						{Suit: model.SuitClub, Number: 1},
					},
				},
				Discards: model.Discards{
					Cards: model.Cards{
						{Suit: model.SuitClub, Number: 10},
					},
				},
			},
		},
		{
			name: "PlayerId not found",
			fields: fields{
				playerId: "10",
				card:     model.Card{},
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
					ActionType: TypePutCard,
					Reason:     "PlayerId is not found.",
				},
			},
			wantState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
			},
		},
		{
			name: "Not turn owner",
			fields: fields{
				playerId: "2",
				card:     model.Card{},
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
					ActionType: TypePutCard,
					Reason:     "PlayerId is not turn owner.",
				},
			},
			wantState: &state.State{
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
			name: "No upcard",
			fields: fields{
				playerId: "2",
				card:     model.Card{},
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 1,
					},
					Upcards: model.Upcards{
						Cards: model.Cards{},
					},
				},
			},
			wantEvents: event.Events{
				&event.NotAcceptable{
					ActionType: TypePutCard,
					Reason:     "No upcards exists.",
				},
			},
			wantState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 1,
				},
				Upcards: model.Upcards{
					Cards: model.Cards{},
				},
			},
		},
		{
			name: "Different number and different suit.",
			fields: fields{
				playerId: "2",
				card:     model.Card{Suit: model.SuitClub, Number: 10},
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 1,
					},
					Upcards: model.Upcards{
						Cards: model.Cards{
							{Suit: model.SuitSpade, Number: 11},
						},
					},
				},
			},
			wantEvents: event.Events{
				&event.NotAcceptable{
					ActionType: TypePutCard,
					Reason:     "PutCard has different number and different suit.",
				},
			},
			wantState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 1,
				},
				Upcards: model.Upcards{
					Cards: model.Cards{
						{Suit: model.SuitSpade, Number: 11},
					},
				},
			},
		},
		{
			name: "Has no PutCard",
			fields: fields{
				playerId: "2",
				card:     model.Card{Suit: model.SuitSpade, Number: 10},
			},
			args: args{
				st: &state.State{
					Players: []model.Player{
						{Id: "1"},
						{Id: "2"},
					},
					PlayOrder: model.PlayOrder{
						Order:      []model.PlayerId{"1", "2"},
						CurrentIdx: 1,
					},
					Upcards: model.Upcards{
						Cards: model.Cards{
							{Suit: model.SuitSpade, Number: 11},
						},
					},
				},
			},
			wantEvents: event.Events{
				&event.NotAcceptable{
					ActionType: TypePutCard,
					Reason:     "Player has not PutCard.",
				},
			},
			wantState: &state.State{
				Players: []model.Player{
					{Id: "1"},
					{Id: "2"},
				},
				PlayOrder: model.PlayOrder{
					Order:      []model.PlayerId{"1", "2"},
					CurrentIdx: 1,
				},
				Upcards: model.Upcards{
					Cards: model.Cards{
						{Suit: model.SuitSpade, Number: 11},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := &PutCard{
				PlayerId: tt.fields.playerId,
				Card:     tt.fields.card,
			}
			gotEvents, gotState := pc.Do(tt.args.st)
			if !reflect.DeepEqual(gotEvents, tt.wantEvents) {
				t.Errorf("PutCard.Do() = %+v, want %+v", gotEvents, tt.wantEvents)
			}
			if !reflect.DeepEqual(gotState, tt.wantState) {
				t.Errorf("PutCard.Do() = %+v, want %+v", gotState, tt.wantState)
			}
		})
	}
}
