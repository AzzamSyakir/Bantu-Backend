package controllers

import (
	"bantu-backend/src/internal/helper/response"
	"bantu-backend/src/internal/services"
)

type UserController struct {
	UserService     *services.UserService
	ResponseChannel chan response.Response[any]
}

func NewUserController(authService *services.UserService) *UserController {
	return &UserController{
		UserService:     authService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}
