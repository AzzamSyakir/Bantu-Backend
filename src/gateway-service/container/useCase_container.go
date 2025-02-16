package container

import (
	"bantu-backend/src/gateway-service/use_case"
)

type UseCaseContainer struct {
	Gateway *use_case.GatewayUseCase
	Expose  *use_case.ExposeUseCase
}

func NewUseCaseContainer(
	auth *use_case.GatewayUseCase,
	expose *use_case.ExposeUseCase,

) *UseCaseContainer {
	useCaseContainer := &UseCaseContainer{
		Gateway: auth,
		Expose:  expose,
	}
	return useCaseContainer
}
