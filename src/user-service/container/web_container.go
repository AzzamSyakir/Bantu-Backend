package container

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/user-service/config"
	"bantu-backend/src/user-service/delivery/grpc/client"
	"bantu-backend/src/user-service/repository"
	"bantu-backend/src/user-service/use_case"
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
		envConfig.App.AuthHost,
		envConfig.App.AuthGrpcPort,
	)
	initAuthClient := client.InitAuthServiceClient(authUrl)
	userUseCase := use_case.NewUserUseCase(&initAuthClient, userDBConfig, userRepository)
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
