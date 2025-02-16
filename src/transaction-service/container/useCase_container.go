package container

import (
	"bantu-backend/src/transaction-service/use_case"
)

type UseCaseContainer struct {
	product  *use_case.ProductUseCase
	Category *use_case.CategoryUseCase
}

func NewUseCaseContainer(
	product *use_case.ProductUseCase,
	category *use_case.CategoryUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		product:  product,
		Category: category,
	}
	return useCaseContainer
}
