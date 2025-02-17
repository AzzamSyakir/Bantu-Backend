package controllers

import "bantu-backend/src/internal/services"

type UserController struct {
	UserService *services.UserService
}

func NewUserController(authService *services.UserService) *UserController {
	return &UserController{
		UserService: authService,
	}
}
