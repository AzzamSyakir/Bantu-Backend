package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type UserService struct {
	UserRepository *repository.UserRepository
	Producer       *producer.ServicesProducer
}

func NewUserService(userRepository *repository.UserRepository, producer *producer.ServicesProducer) *UserService {
	UserService := &UserService{
		Producer:       producer,
		UserRepository: userRepository,
	}
	return UserService
}
