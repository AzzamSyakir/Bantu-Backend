package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type AuthService struct {
	UserRepository *repository.UserRepository
	Producer       *producer.ServicesProducer
}

func NewAuthService(userRepository *repository.UserRepository, producer *producer.ServicesProducer) *AuthService {
	AuthService := &AuthService{
		Producer:       producer,
		UserRepository: userRepository,
	}
	return AuthService
}
