package app

import (
	"fmt"
	"net/http"

	"mrps-game/internal/app/logic/model"
	httpTransport "mrps-game/internal/app/transport/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Start(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Migrate the schema
	err = db.AutoMigrate(&model.User{}, &model.Transaction{})
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	server := httpTransport.NewServer(db)
	go server.Listen()

	return http.ListenAndServe(":8080", nil)
}
