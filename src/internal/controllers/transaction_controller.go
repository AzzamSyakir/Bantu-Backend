package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type TransactionController struct {
	TransactionService *services.TransactionService
	ResponseChannel    *response.ResponseChannel
}

func NewTransactionController(jobService *services.TransactionService, responseChannel *response.ResponseChannel) *TransactionController {
	return &TransactionController{
		TransactionService: jobService,
		ResponseChannel:    responseChannel,
	}
}
