package container

import (
	"bantu-backend/src/gateway-service/delivery/http"
)

type ControllerContainer struct {
	Auth   *http.AuthController
	Expose *http.ExposeController
}

func NewControllerContainer(
	auth *http.AuthController,
	expose *http.ExposeController,

) *ControllerContainer {
	controllerContainer := &ControllerContainer{
		Auth:   auth,
		Expose: expose,
	}
	return controllerContainer
}
