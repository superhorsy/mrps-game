package logic

import (
	"encoding/json"
	"fmt"

	"mrps-game/internal/app/logic/message"
)

type GameError struct {
	errorText string
}

func (g GameError) Error() string {
	return fmt.Sprintf("game error: %s", g.errorText)
}

func (g GameError) ToResponse() []byte {
	errResponseJSON, _ := json.Marshal(message.NewError(g.errorText))
	return errResponseJSON
}

func NewGameError(errorText string) GameError {
	return GameError{
		errorText: errorText,
	}
}
