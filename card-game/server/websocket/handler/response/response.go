package response

type Type string

const (
	TypeRoomInfo Type = "RoomInfo"
)

type Response struct {
	Type Type
	Body interface{}
}
