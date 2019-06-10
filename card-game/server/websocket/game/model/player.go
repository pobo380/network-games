package model

type PlayerId string

type Player struct {
	Id PlayerId

	Hand Hand
}
