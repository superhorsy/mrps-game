package http

import (
	"log"
	"net/http"

	"mrps-game/internal/app/logic"
	"mrps-game/internal/app/service"
	"mrps-game/internal/app/transport/http/handler"

	"github.com/rs/cors"
	"gorm.io/gorm"
)

// Server struct holds server variables.
type Server struct {
	db *gorm.DB
}

// NewServer initializes a new server.
func NewServer(db *gorm.DB) *Server {
	return &Server{
		db: db,
	}
}

// Listen runs puts server into listening mode.
func (s *Server) Listen() {
	log.Println("Start HTTP server...")

	userService := service.NewUserService(s.db)
	transService := service.NewTransactionService(s.db)
	userHandler := handler.NewUserHandler(userService)
	gameServer := logic.NewServer(transService)

	websocketHandler := handler.NewWebsocketHandler(userService, gameServer)

	http.HandleFunc("/ws", websocketHandler.Handle)
	http.Handle("/register", cors.Default().Handler(http.HandlerFunc(userHandler.Register)))
	http.Handle("/login", cors.Default().Handler(http.HandlerFunc(userHandler.Login)))
}
