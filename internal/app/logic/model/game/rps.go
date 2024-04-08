package game

// RPS is a Rock-Paper-Scissors game.

type RPSGame struct {
	Move1 string
	Move2 string
}

func NewRPSGame() *RPSGame {
	return &RPSGame{}
}

func (r *RPSGame) Name() string {
	return "rps"
}

func (r *RPSGame) SetMoveForPlayer1(move interface{}) error {
	if moveStr, ok := move.(string); ok {
		r.Move1 = moveStr
		return nil
	}
	return ErrInvalidMove
}

func (r *RPSGame) SetMoveForPlayer2(move interface{}) error {
	if moveStr, ok := move.(string); ok {
		r.Move2 = moveStr
		return nil
	}
	return ErrInvalidMove
}

func (r *RPSGame) Result() int {
	if r.Move1 == r.Move2 {
		return 0
	}
	if r.Move1 == "rock" && r.Move2 == "scissors" {
		return 1
	}
	if r.Move1 == "scissors" && r.Move2 == "paper" {
		return 1
	}
	if r.Move1 == "paper" && r.Move2 == "rock" {
		return 1
	}
	return -1
}
