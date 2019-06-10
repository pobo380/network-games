package event

const (
	TypeNotAcceptable = "NotAcceptable"

	TypeTurnChange = "TurnChange"
	TypePutCard    = "PutCard"
	TypePass       = "Pass"
	TypeDrawCard   = "DrawCard"
	TypeFinishGame = "FinishGame"
)

type Event interface{}

type Events []Event

func (es *Events) Add(e Event) {
	*es = append(*es, e)
}
