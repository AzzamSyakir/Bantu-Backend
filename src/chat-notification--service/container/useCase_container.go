package container

import (
	"bantu-backend/src/chat-notification-service/use_case"
)

type UseCaseContainer struct {
	User *use_case.UserUseCase
}

func NewUseCaseContainer(
	user *use_case.UserUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		User: user,
	}
	return useCaseContainer
}
