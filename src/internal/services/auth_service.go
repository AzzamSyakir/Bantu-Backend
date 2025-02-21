package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"errors"
	"regexp"
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

	if request.Email == "" || request.Name == "" || request.Password == "" {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("email, name and password cant be empty")
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("invalid email")
	}

	if request.Role == "" {
		request.Role = "client"
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, hashedPasswordErr
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

func (authService *AuthService) LoginService(request *request.LoginRequest) (result *entity.UserEntity, err error) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		return nil, err
		// authService.Producer.CreateMessageAuth(authService.EnvConfig.RabbitMq, err.Error())
	}

	if request.Email == "" || request.Password == "" {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("email and password cant be empty")
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("invalid email")
	}

	foundUser, err := authService.UserRepository.LoginUser(begin, request.Email)
	if err != nil {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}

	if foundUser.Email == "" {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(request.Password))
	if err != nil {
		rollbackErr := begin.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, errors.New("invalid password")
	}

	return foundUser, nil
}

// func (authService *AuthService) GenerateToken(email string) (tokenString string, err error) {

// 	claims := jwt.MapClaims{
// 		"em":  email,
// 		"exp": time.Now().Add(time.Hour * 1).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	tokenString, err := token.SignedString(jwtSecret)
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func (authService *AuthService) ValidateToken(token string) (token string, err error) {

// }
