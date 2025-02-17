package services

import "bantu-backend/src/internal/repository"

type AuthService struct {
	UserRepository *repository.UserRepository
}

func NewAuthService(userRepository *repository.UserRepository) *AuthService {
	AuthService := &AuthService{
		UserRepository: userRepository,
	}
	return AuthService
}
