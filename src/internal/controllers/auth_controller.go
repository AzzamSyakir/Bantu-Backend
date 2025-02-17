package controllers

import "bantu-backend/src/internal/services"

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}
