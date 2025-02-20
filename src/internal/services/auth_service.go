package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DatabaseConfig *configs.DatabaseConfig
	Rabbitmq       *configs.RabbitMqConfig
}

func NewAuthService(userRepository *repository.UserRepository, producer *producer.ServicesProducer, envConfig *configs.EnvConfig, dbConfig *configs.DatabaseConfig) *AuthService {
	AuthService := &AuthService{
		DatabaseConfig: dbConfig,
		EnvConfig:      envConfig,
		Producer:       producer,
		UserRepository: userRepository,
		Rabbitmq:       rabbitmq,
	}
	return AuthService
}

func (authService *AuthService) RegisterService(request *request.RegisterRequest) (result *entity.UserEntity, err error) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		return nil, err
		// authService.Producer.CreateMessageAuth(authService.EnvConfig.RabbitMq, err.Error())
	}

	if request.Email.IsZero() || request.Name.IsZero() || request.Password.IsZero() {
		rollbackErr := begin.Rollback()
		return nil, rollbackErr
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		rollbackErr := begin.Rollback()
		return nil, rollbackErr
	}

	currentTime := null.NewTime(time.Now(), true)
	newUser := &entity.UserEntity{
		ID:        null.NewString(string(uuid.NewString()), true),
		Name:      request.Name,
		Email:     request.Email,
		Password:  null.NewString(string(hashedPassword), true),
		Role:      request.Role,
		Balance:   null.NewFloat(0.0, true),
		CreatedAt: currentTime.Time,
		UpdatedAt: currentTime.Time,
	}

	createdUser, err := authService.UserRepository.RegisterUser(begin, newUser)
	if err != nil {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		return nil, commitErr
	}
	return createdUser, nil
}
