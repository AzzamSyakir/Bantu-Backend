package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type UserController struct {
	UserService     *services.UserService
	ResponseChannel *response.ResponseChannel
}

func NewUserController(authService *services.UserService, responseChannel *response.ResponseChannel) *UserController {
	return &UserController{
		UserService:     authService,
		ResponseChannel: responseChannel,
	}
}
