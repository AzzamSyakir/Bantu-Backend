package consumer

import (
	"bantu-backend/src/internal/services"
)

type ServiceConsumer struct {
	AuthService        *services.AuthService
	ChatService        *services.ChatService
	JobService         *services.JobService
	ProposalService    *services.ProposalService
	TransactionService *services.TransactionService
	UserService        *services.UserService
}

func NewServiceConsumer(
	authService *services.AuthService,
	chatService *services.ChatService,
	jobService *services.JobService,
	proposalService *services.ProposalService,
	transactionService *services.TransactionService,
	userService *services.UserService,
) *ServiceConsumer {
	return &ServiceConsumer{
		AuthService:        authService,
		JobService:         jobService,
		ProposalService:    proposalService,
		TransactionService: transactionService,
		UserService:        userService,
	}
}
