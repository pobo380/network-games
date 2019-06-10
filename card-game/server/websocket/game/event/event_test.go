package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvents_Add(t *testing.T) {
	type args struct {
		e Event
	}
	tests := []struct {
		name string
		es   Events
		args args
	}{
		{
			name: "正常系",
			es: Events{
				&NotAcceptable{},
			},
			args: args{
				e: &NotAcceptable{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.es.Add(tt.args.e)

			assert.Len(t, tt.es, 2)
		})
	}
}
