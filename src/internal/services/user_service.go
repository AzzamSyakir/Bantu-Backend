package services

import "bantu-backend/src/internal/repository"

type UserService struct {
	UserRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	UserService := &UserService{
		UserRepository: userRepository,
	}
	return UserService
}
