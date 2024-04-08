package message

import (
	"time"

	"mrps-game/internal/app/logic/model"
)

type FundsAdded struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
	Total  int64  `json:"total"`
}

func NewFundsAdded(amount, total int64) *FundsAdded {
	return &FundsAdded{
		Type:   "funds.added",
		Amount: amount,
		Total:  total,
	}
}

type FundsWithdrawn struct {
	Type   string `json:"type"`
	Amount int64  `json:"amount"`
	Total  int64  `json:"total"`
}

func NewFundsWithdrawn(amount, total int64) *FundsWithdrawn {
	return &FundsWithdrawn{
		Type:   "funds.withdrawn",
		Amount: amount,
		Total:  total,
	}
}

type Error struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

func NewError(err string) *Error {
	return &Error{Error: err, Type: "error"}
}

type ChallengeList []Challenge

type Challenge struct {
	ChallengerId uint32 `json:"challengerId"`
	OpponentId   uint32 `json:"opponentId"`
	Amount       int64  `json:"amount"`
	Game         string `json:"game"`
}

type ClientInfo struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Funds int64  `json:"funds"`
	Id    uint32 `json:"id"`
}

func NewClientInfo(name string, funds int64, id uint32) *ClientInfo {
	return &ClientInfo{
		Type:  "me",
		Name:  name,
		Funds: funds,
		Id:    id,
	}
}

type Transaction struct {
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
type Opponents struct {
	Type      string           `json:"type"`
	Opponents []model.Opponent `json:"opponents"`
}

func NewOpponents(opponents []model.Opponent) *Opponents {
	return &Opponents{
		Type:      "opponents",
		Opponents: opponents,
	}
}

type GameResult struct {
	Type   string `json:"type"`
	Result string `json:"result"`
	Amount int64  `json:"amount"`
}

func NewGameResult(result string, amount int64) *GameResult {
	return &GameResult{
		Type:   "game.result",
		Result: result,
		Amount: amount,
	}
}
