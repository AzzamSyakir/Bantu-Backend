package container

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/chat-notification-service/config"
	"bantu-backend/src/chat-notification-service/delivery/grpc/client"
	"bantu-backend/src/chat-notification-service/repository"
	"bantu-backend/src/chat-notification-service/use_case"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env        *config.EnvConfig
	UserDB     *config.DatabaseConfig
	Repository *RepositoryContainer
	UseCase    *UseCaseContainer
	Grpc       *grpc.Server
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	userDBConfig := config.NewUserDBConfig(envConfig)

	userRepository := repository.NewUserRepository()
	repositoryContainer := NewRepositoryContainer(userRepository)
	authUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.GatewayHost,
		envConfig.App.GatewayGrpcPort,
	)
	initGatewayClient := client.InitGatewayServiceClient(authUrl)
	userUseCase := use_case.NewUserUseCase(&initGatewayClient, userDBConfig, userRepository)
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userUseCase)

	useCaseContainer := NewUseCaseContainer(userUseCase)
	webContainer := &WebContainer{
		Env:        envConfig,
		UserDB:     userDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Grpc:       grpcServer,
	}

	return webContainer
}
