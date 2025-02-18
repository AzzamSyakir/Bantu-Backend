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
