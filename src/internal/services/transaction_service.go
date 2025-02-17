package services

import "bantu-backend/src/internal/repository"

type TransactionService struct {
	TransactionRepository *repository.TransactionRepository
}

func NewTransactionService(userRepository *repository.TransactionRepository) *TransactionService {
	TransactionService := &TransactionService{
		TransactionRepository: userRepository,
	}
	return TransactionService
}
