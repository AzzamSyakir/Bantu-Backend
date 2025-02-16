package container

import (
	"bantu-backend/src/review-service/use_case"
)

type UseCaseContainer struct {
	Order *use_case.OrderUseCase
}

func NewUseCaseContainer(
	order *use_case.OrderUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Order: order,
	}
	return useCaseContainer
}
