package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type TransactionController struct {
	TransactionService *services.TransactionService
	ResponseChannel    chan response.Response[any]
}

func NewTransactionController(jobService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: jobService,
		ResponseChannel:    make(chan response.Response[any], 1),
	}
}
