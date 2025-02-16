package container

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/review-service/config"
	"bantu-backend/src/review-service/delivery/grpc/client"
	"bantu-backend/src/review-service/repository"
	"bantu-backend/src/review-service/use_case"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env        *config.EnvConfig
	OrderDB    *config.DatabaseConfig
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
	orderDBConfig := config.NewDBConfig(envConfig)

	orderRepository := repository.NewOrderRepository()
	repositoryContainer := NewRepositoryContainer(orderRepository)
	userUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.UserHost,
		envConfig.App.UserPort,
	)
	productUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.ProductHost,
		envConfig.App.ProductPort,
	)
	initUserClient := client.InitUserServiceClient(userUrl)
	initProductClient := client.InitProductServiceClient(productUrl)
	orderUseCase := use_case.NewOrderUseCase(orderDBConfig, orderRepository, envConfig, &initUserClient, &initProductClient)

	useCaseContainer := NewUseCaseContainer(orderUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, orderUseCase)

	webContainer := &WebContainer{
		Env:        envConfig,
		OrderDB:    orderDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Grpc:       grpcServer,
	}

	return webContainer
}
