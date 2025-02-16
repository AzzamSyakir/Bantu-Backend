package container

import (
	"bantu-backend/src/gateway-service/delivery/http"
)

type ControllerContainer struct {
	Gateway *http.GatewayController
	Expose  *http.ExposeController
}

func NewControllerContainer(
	auth *http.GatewayController,
	expose *http.ExposeController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Gateway: auth,
		Expose:  expose,
	}
	return controllerContainer
}
