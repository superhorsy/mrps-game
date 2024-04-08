package message

type ChallengeSend struct {
	Amount   float64 `json:"amount"`
	Opponent float64 `json:"opponent"`
	Move     string  `json:"move"`
	Game     string  `json:"game"`
}
