package model

type PlayOrder struct {
	Order      []PlayerId
	CurrentIdx int
}

func (po *PlayOrder) CurrentPlayerId() PlayerId {
	return po.Order[po.CurrentIdx]
}

func (po *PlayOrder) NextPlayerId() PlayerId {
	return po.Order[po.nextIdx()]
}

func (po *PlayOrder) StepToNextPlayer() PlayerId {
	po.CurrentIdx = po.nextIdx()
	return po.CurrentPlayerId()
}

func (po *PlayOrder) nextIdx() int {
	return (po.CurrentIdx + 1) % len(po.Order)
}
