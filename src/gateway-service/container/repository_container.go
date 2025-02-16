package container

import (
	"bantu-backend/src/gateway-service/repository"
)

type RepositoryContainer struct {
	Auth *repository.AuthRepository
}

func NewRepositoryContainer(
	auth *repository.AuthRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Auth: auth,
	}
	return repositoryContainer
}
