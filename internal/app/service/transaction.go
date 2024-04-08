package service

import (
	"time"

	"mrps-game/internal/app/logic/model"

	"gorm.io/gorm"
)

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) LogTransaction(userID uint32, amount float64, reason string) error {
	transactionType := "add"
	if amount < 0 {
		transactionType = "subtract"
		amount = -amount // Ensure the amount is positive for storage
	}

	transaction := model.Transaction{
		UserID:    userID,
		Type:      transactionType,
		Reason:    reason,
		Amount:    amount,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&transaction).Error; err != nil {
		return err
	}

	return nil
}
func (s *TransactionService) GetLastTransactions(userID uint32) ([]model.Transaction, error) {
	var transactions []model.Transaction

	err := s.db.Where("user_id = ?", userID).Order("created_at desc").Limit(10).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
