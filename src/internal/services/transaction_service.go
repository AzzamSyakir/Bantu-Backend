package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type TransactionService struct {
	TransactionRepository *repository.TransactionRepository
	Producer              *producer.ServicesProducer
}

func NewTransactionService(userRepository *repository.TransactionRepository, producer *producer.ServicesProducer) *TransactionService {
	TransactionService := &TransactionService{
		Producer:              producer,
		TransactionRepository: userRepository,
	}
	return TransactionService
}
