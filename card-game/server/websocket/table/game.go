package table

type GameId string

type Game struct {
	GameId    string
	PlayerIds []string

	RawState string
}
