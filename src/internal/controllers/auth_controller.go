package controllers

import (
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/services"
	"net/http"
)

type AuthController struct {
	AuthService     *services.AuthService
	ResponseChannel chan response.Response[any]
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService:     authService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}

func (authController *AuthController) Register(writer http.ResponseWriter, reader *http.Request) {

	request := request.RegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	service := authController.AuthService.RegisterService(request)

	response.NewResponse[any](writer, &response.Response[any]{
		Code:    http.StatusOK,
		Message: "Register success",
		Data:    nil,
	})
}
