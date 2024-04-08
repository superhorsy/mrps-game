package handler

import (
	"log"
	"net/http"

	"mrps-game/internal/app/logic"
	"mrps-game/internal/app/service"
	"mrps-game/internal/app/utils"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
)

type WebsocketHandler struct {
	userService *service.UserService
	upgrader    *websocket.Upgrader
	server      *logic.GameServer
}

func NewWebsocketHandler(svc *service.UserService, server *logic.GameServer) *WebsocketHandler {
	return &WebsocketHandler{
		server:      server,
		userService: svc,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			}},
	}
}

func (h *WebsocketHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the URL.
	token := r.URL.Query().Get("token")

	// Verify the token.
	if !utils.VerifyToken(token) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Read token claims.
	claims, err := utils.ReadToken(token)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Get the user ID from the claims.
	userIDUint32 := cast.ToUint32(claims["user_id"])

	// Get user from the database.
	user, err := h.userService.GetUserByID(userIDUint32)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := logic.NewClient(uint32(user.ID), user.Username, conn, h.server)
	h.server.Clients.Add(client)

	log.Println("Added new client. Now", h.server.Clients.Count(), "clients connected.")
	client.Listen()
}
