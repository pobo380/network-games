package request

import (
	"encoding/json"
)

type Request struct {
	Action string `json:"action"`
	Body   interface{}
}

func Parse(bts []byte, body interface{}) error {
	req := &Request{Body: body}
	err := json.Unmarshal(bts, req)
	if err != nil {
		return err
	}
	return nil
}
