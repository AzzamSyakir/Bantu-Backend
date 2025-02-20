package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DatabaseConfig  *configs.DatabaseConfig
	Rabbitmq        *configs.RabbitMqConfig
	EnvConfig       *configs.EnvConfig
	Producer        *producer.ServicesProducer
	UserRepository  *repository.UserRepository
	AdminRepository *repository.AdminRepository
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
		errMessage := "email, name and password must be provided"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		begin.Rollback()
		errMessage := "invalid email type"
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	if request.Role == "" {
		request.Role = "client"
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, hashedPasswordErr.Error(), http.StatusInternalServerError)
		return
	}

	newUser := &entity.UserEntity{
		ID:        string(uuid.NewString()),
		Name:      request.Name,
		Email:     request.Email,
		Password:  string(hashedPassword),
		Role:      request.Role,
		Balance:   0.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, createdUserErr := authService.UserRepository.RegisterUser(begin, newUser)
	if createdUserErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, createdUserErr.Error(), http.StatusInternalServerError)
		return
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}

	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, createdUser)

}

func (authService *AuthService) LoginService(request *request.LoginRequest) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Email == "" || request.Password == "" {
		errMessage := "email and password must be provided"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		errMessage := "invalid email type"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	foundUser, foundUserErr := authService.UserRepository.LoginUser(begin, request.Email)
	if foundUserErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, foundUserErr.Error(), http.StatusInternalServerError)
		return
	}

	if foundUser == nil {
		errMessage := "user not found"
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(request.Password))
	if err != nil {
		errMessage := "invalid password"
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusUnauthorized)
		return
	}

	tokenString, err := authService.GenerateToken(foundUser.ID, foundUser.Role)
	if err != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return
	}

	foundUser.Token = tokenString
	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, foundUser)
}

func (authService *AuthService) GenerateToken(id string, role string) (string, error) {
	jwtSecret := []byte(authService.EnvConfig.SecretKey)
	claims := jwt.MapClaims{
		"rl":  role,
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (authService *AuthService) AdminRegisterService(request *request.AdminRegisterRequest) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
	}

	if request.Email == "" || request.Password == "" || request.Username == "" {
		errMessage := "email, username and password must be provided"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		errMessage := "invalid email type"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, hashedPasswordErr.Error(), http.StatusInternalServerError)
		return
	}

	currentTime := time.Now()
	newAdmin := &entity.AdminEntity{
		ID:        string(uuid.NewString()),
		Username:  request.Username,
		Email:     request.Email,
		Password:  string(hashedPassword),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	createdAdmin, createdAdminErr := authService.AdminRepository.RegisterAdmin(begin, newAdmin)
	if createdAdminErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, createdAdminErr.Error(), http.StatusInternalServerError)
		return
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}

	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, createdAdmin)

}

func (authService *AuthService) AdminLoginService(request *request.AdminLoginRequest) {
	begin, err := authService.DatabaseConfig.DB.Connection.Begin()
	if err != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.Email == "" || request.Password == "" {
		errMessage := "email and password must be provided"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	emailRegex := `^(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(request.Email) {
		errMessage := "invalid email type"
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}

	foundAdmin, foundAdminErr := authService.AdminRepository.LoginAdmin(begin, request.Email)
	if foundAdminErr != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, foundAdminErr.Error(), http.StatusInternalServerError)
		return
	}

	if foundAdmin == nil {
		errMessage := "user not found"
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundAdmin.Password), []byte(request.Password))
	if err != nil {
		errMessage := "invalid password"
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, errMessage, http.StatusUnauthorized)
		return
	}

	tokenString, err := authService.GenerateToken(foundAdmin.ID, "admin")
	if err != nil {
		begin.Rollback()
		authService.Producer.CreateMessageError(authService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return
	}

	foundAdmin.Token = tokenString
	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, foundAdmin)

}
