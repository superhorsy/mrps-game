package game

import (
	"errors"
)

// DiceGame is a random game, where the player with the highest number wins.
// Number are taken randomly from 1 to 6 and compared.
type DiceGame struct {
	Move1 int
	Move2 int
}

var ErrInvalidMove = errors.New("invalid move")

func NewDiceGame() *DiceGame {
	return &DiceGame{}
}

func (d *DiceGame) Name() string {
	return "dice"
}

func (d *DiceGame) SetMoveForPlayer1(move interface{}) error {
	if moveInt, ok := move.(int); ok {
		d.Move1 = moveInt
		return nil
	}
	return ErrInvalidMove
}

func (d *DiceGame) SetMoveForPlayer2(move interface{}) error {
	if moveInt, ok := move.(int); ok {
		d.Move2 = moveInt
		return nil
	}
	return ErrInvalidMove
}

func (d *DiceGame) Result() int {
	if d.Move1 == d.Move2 {
		return 0
	}
	if d.Move1 > d.Move2 {
		return 1
	}
	return -1
}
