package model

type Challenge struct {
	ChallengerId uint32
	OpponentId   uint32
	Amount       int64
	Game         Game
	Move         interface{}
}

type Game interface {
	Name() string
	Result() int
	SetMoveForPlayer1(move interface{}) error
	SetMoveForPlayer2(move interface{}) error
}
