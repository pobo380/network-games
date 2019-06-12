package table

type Room struct {
	RoomId string
	IsOpen string `dynamodbav:",omitempty"`

	MaxPlayerNum int
	PlayerIds    []string
}

func (r *Room) AddPlayer(playerId string) {
	r.PlayerIds = append(r.PlayerIds, playerId)

	if l := len(r.PlayerIds); l >= r.MaxPlayerNum {
		r.Close()
	}
}

func (r *Room) Close() {
	r.IsOpen = ""
}

func (r *Room) IsClosed() bool {
	return r.IsOpen == ""
}
