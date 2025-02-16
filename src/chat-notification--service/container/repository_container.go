package container

import (
	"bantu-backend/src/chat-notification-service/repository"
)

type RepositoryContainer struct {
	User *repository.UserRepository
}

func NewRepositoryContainer(
	user *repository.UserRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		User: user,
	}
	return repositoryContainer
}
