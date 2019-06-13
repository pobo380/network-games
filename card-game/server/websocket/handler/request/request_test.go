package request

import (
	"reflect"
	"testing"
)

type testBodyStruct struct {
	A int
}

func Test_Parse(t *testing.T) {
	type args struct {
		bts  []byte
		body interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantBody interface{}
		wantErr  bool
	}{
		{
			name: "正常系",
			args: args{
				bts:  []byte(`{"Body": {"A": 1}}`),
				body: &testBodyStruct{},
			},
			wantBody: &testBodyStruct{A: 1},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Parse(tt.args.bts, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.body, tt.wantBody) {
				t.Errorf("Parse() body = %+v, wantBody %+v", tt.args.body, tt.wantBody)
			}
		})
	}
}
