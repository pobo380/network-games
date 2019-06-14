package event

const (
	TypeNotAcceptable = "NotAcceptable"

	TypeTurnChange = "TurnChange"
	TypePutCard    = "PutCard"
	TypePass       = "Pass"
	TypeDrawCard   = "DrawCard"
	TypeGameState  = "GameState"
	TypeFinishGame = "FinishGame"
)

type Event interface {
	GetType() string
}

type Events []Event

func (es *Events) Add(e Event) {
	*es = append(*es, e)
}
