package controllers

import "bantu-backend/src/internal/services"

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(jobService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: jobService,
	}
}
