package container

import (
	"bantu-backend/src/internal/controllers"
)

type ControllerContainer struct {
	AuthController        *controllers.AuthController
	ChatController        *controllers.ChatController
	JobController         *controllers.JobController
	ProposalController    *controllers.ProposalController
	TransactionController *controllers.TransactionController
	UserController        *controllers.UserController
}

func NewControllerContainer(
	authController *controllers.AuthController,
	userController *controllers.UserController,
	chatController *controllers.ChatController,
	jobController *controllers.JobController,
	proposalController *controllers.ProposalController,
	transactionController *controllers.TransactionController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		AuthController:        authController,
		ChatController:        chatController,
		JobController:         jobController,
		ProposalController:    proposalController,
		TransactionController: transactionController,
		UserController:        userController,
	}
	return controllerContainer
}
