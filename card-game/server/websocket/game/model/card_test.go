package model

import (
	"reflect"
	"testing"
)

func TestCards_Remove(t *testing.T) {
	type args struct {
		target Card
	}
	tests := []struct {
		name      string
		cs        *Cards
		args      args
		want      bool
		wantCards *Cards
	}{
		{
			name: "正常系",
			cs: &Cards{
				{Suit: SuitClub, Number: 1},
				{Suit: SuitSpade, Number: 2},
			},
			args: args{
				target: Card{Suit: SuitClub, Number: 1},
			},
			want: true,
			wantCards: &Cards{
				{Suit: SuitSpade, Number: 2},
			},
		},
		{
			name: "みつからない",
			cs: &Cards{
				{Suit: SuitClub, Number: 1},
				{Suit: SuitSpade, Number: 2},
			},
			args: args{
				target: Card{Suit: SuitClub, Number: 11},
			},
			want: false,
			wantCards: &Cards{
				{Suit: SuitClub, Number: 1},
				{Suit: SuitSpade, Number: 2},
			},
		},
		{
			name: "len(cs) == 0",
			cs:   &Cards{},
			args: args{
				target: Card{Suit: SuitClub, Number: 11},
			},
			want:      false,
			wantCards: &Cards{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cs.Remove(tt.args.target); got != tt.want {
				t.Errorf("Cards.Remove() = %+v, want %+v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.cs, tt.wantCards) {
				t.Errorf("Cards.Remove() = %+v, want %+v", tt.cs, tt.wantCards)
			}
		})
	}
}

func TestCards_RemoveTop(t *testing.T) {
	tests := []struct {
		name      string
		cs        *Cards
		wantCards *Cards
	}{
		{
			name: "正常系",
			cs: &Cards{
				{Suit: SuitClub, Number: 1},
				{Suit: SuitSpade, Number: 2},
			},
			wantCards: &Cards{
				{Suit: SuitClub, Number: 1},
			},
		},
		{
			name:      "手札がない",
			cs:        &Cards{},
			wantCards: &Cards{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cs.RemoveTop()
			if !reflect.DeepEqual(tt.cs, tt.wantCards) {
				t.Errorf("Cards.Remove() = %+v, want %+v", tt.cs, tt.wantCards)
			}
		})
	}
}
