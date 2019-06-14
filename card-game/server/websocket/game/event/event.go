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

func NewFromType(t string) Event {
	switch t {
	case TypeNotAcceptable:
		return &NotAcceptable{}
	case TypeTurnChange:
		return &TurnChange{}
	case TypePutCard:
		return &PutCard{}
	case TypePass:
		return &Pass{}
	case TypeDrawCard:
		return &DrawCard{}
	case TypeGameState:
		return &GameState{}
	case TypeFinishGame:
		return &FinishGame{}
	}

	return nil
}

type Events []Event

func (es *Events) Add(e Event) {
	*es = append(*es, e)
}
