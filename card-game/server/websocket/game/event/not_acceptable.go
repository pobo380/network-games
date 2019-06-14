package event

type NotAcceptable struct {
	ActionType string
	Reason     string
}

func (*NotAcceptable) GetType() string {
	return TypeNotAcceptable
}
