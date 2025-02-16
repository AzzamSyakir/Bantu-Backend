package container

import (
	"bantu-backend/src/review-service/repository"
)

type RepositoryContainer struct {
	Order *repository.OrderRepository
}

func NewRepositoryContainer(
	order *repository.OrderRepository,

) *RepositoryContainer {
	repositoryContainer := &RepositoryContainer{
		Order: order,
	}
	return repositoryContainer
}
