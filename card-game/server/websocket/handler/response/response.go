package response

type Type string

const (
	TypeRoomInfo  Type = "RoomInfo"
	TypeGameStart Type = "GameStart"
	TypeGameEvent Type = "GameEvent"
)

type Response struct {
	Type Type
	Body interface{}
}

type Responses []*Response

func (r *Responses) Add(t Type, body interface{}) {
	*r = append(*r, &Response{
		Type: t,
		Body: body,
	})
}
