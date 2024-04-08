package logic

import (
	"log"
	"net/http"

	"mrps-game/internal/app/service"

	"github.com/gorilla/websocket"
)

// GameServer struct holds server variables.
type GameServer struct {
	Clients      *Clients
	upgrader     *websocket.Upgrader
	transService *service.TransactionService
}

// NewServer initializes a new server.
func NewServer(svc *service.TransactionService) *GameServer {
	return &GameServer{
		transService: svc,
		Clients:      NewClients(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
	}
}

func (s *GameServer) SendToClient(clientID uint32, message string) {
	client, ok := s.Clients.Get(clientID)
	if ok {
		client.sendMessage([]byte(message))
	} else {
		log.Printf("Client %d not found\n", clientID)
		return
	}
}
