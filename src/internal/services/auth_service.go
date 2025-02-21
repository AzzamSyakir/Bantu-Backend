package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DatabaseConfig *configs.DatabaseConfig
	Rabbitmq       *configs.RabbitMqConfig
	EnvConfig      *configs.EnvConfig
	Producer       *producer.ServicesProducer
	UserRepository *repository.UserRepository
}

func NewAuthService(userRepository *repository.UserRepository, producer *producer.ServicesProducer, envConfig *configs.EnvConfig, dbConfig *configs.DatabaseConfig, rabbitmq *configs.RabbitMqConfig) *AuthService {
	AuthService := &AuthService{
		UserRepository: userRepository,
		Producer:       producer,
		EnvConfig:      envConfig,
		DatabaseConfig: dbConfig,
		Rabbitmq:       rabbitmq,
	}
	return AuthService
}

func (authService *AuthService) RegisterService(request *request.RegisterRequest) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
	}
	if request.Email == "" || request.Name == "" || request.Password == "" {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return
	}
	currentTime := null.NewTime(time.Now(), true)
	newUser := &entity.UserEntity{
		ID:        string(uuid.NewString()),
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(hashedPassword),
		Role:      request.Role,
		Balance:   0.0,
		CreatedAt: currentTime.Time,
		UpdatedAt: currentTime.Time,
	}
	createdUser, err := authService.UserRepository.RegisterUser(begin, newUser)
	if err != nil {
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}
	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, createdUser)
	return
}
