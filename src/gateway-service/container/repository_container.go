package container

import (
	"bantu-backend/src/gateway-service/repository"
)

type RepositoryContainer struct {
	Gateway *repository.GatewayRepository
}

func NewRepositoryContainer(
	auth *repository.GatewayRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Gateway: auth,
	}
	return repositoryContainer
}
