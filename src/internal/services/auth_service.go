package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type AuthService struct {
	Database       *configs.DatabaseConfig
	EnvConfig      *configs.EnvConfig
	UserRepository *repository.UserRepository
	Producer       *producer.ServicesProducer
	Rabbitmq       *configs.RabbitMqConfig
}

func NewAuthService(userRepository *repository.UserRepository, producer *producer.ServicesProducer, envConfig *configs.EnvConfig, database *configs.DatabaseConfig, rabbitmq *configs.RabbitMqConfig) *AuthService {
	AuthService := &AuthService{
		Database:       database,
		EnvConfig:      envConfig,
		Producer:       producer,
		UserRepository: userRepository,
		Rabbitmq:       rabbitmq,
	}
	return AuthService
}

func (authService *AuthService) RegisterService(request *request.RegisterRequest) {
	// begin, err := authService.Database.DB.Connection.Begin()
	// if err != nil {
	// 	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, err.Error())
	// }

	// if request.Name.IsZero() || request.Email.IsZero() || request.Password.IsZero() {
	// 	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, err.Error())
	// }

	// hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
	// if hashedPasswordErr != nil {
	// 	authService.Producer.CreateMessageAuth(authService.Rabbitmq.Channel, err.Error())
	// }
	// currentTime := null.NewTime(time.Now(), true)
	// newUser := &pb.User{
	// 	Id:        uuid.NewString(),
	// 	Name:      request.Name,
	// 	Email:     request.Email,
	// 	Password:  string(hashedPassword),
	// 	Balance:   request.Balance,
	// 	Role:      request.Role,
	// 	CreatedAt: timestamppb.New(currentTime.Time),
	// 	UpdatedAt: timestamppb.New(currentTime.Time),
	// }

	// createdUser, err := userUseCase.UserRepository.RegisterUser(begin, newUser)
	// if err != nil {
	// 	rollbackErr := begin.Rollback()
	// 	result = &pb.UserResponse{
	// 		Code:    int64(codes.Internal),
	// 		Message: fmt.Sprintf("Failed to insert new user into database. Error: %v. Rollback status: %v", err, rollbackErr),
	// 		Data:    nil,
	// 	}
	// 	return result, rollbackErr
	// }

	// commitErr := begin.Commit()
	// if commitErr != nil {
	// 	result = &pb.UserResponse{
	// 		Code:    int64(codes.Internal),
	// 		Message: fmt.Sprintf("Failed to commit transaction after user creation. Error: %v", commitErr),
	// 		Data:    nil,
	// 	}
	// 	return result, commitErr
	// }

	// result = &pb.UserResponse{
	// 	Code:    int64(codes.OK),
	// 	Message: "User successfully registered.",
	// 	Data:    createdUser,
	// }
	// return AuthService
	return
}
