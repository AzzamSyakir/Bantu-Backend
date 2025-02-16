package container

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/transaction-service/config"
	"bantu-backend/src/transaction-service/repository"
	"bantu-backend/src/transaction-service/use_case"
	"fmt"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env                *config.EnvConfig
	ProductDB          *config.DatabaseConfig
	ProductRepository  *RepositoryContainer
	CategoryRepository *RepositoryContainer
	UseCase            *UseCaseContainer
	Grpc               *grpc.Server
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	productDBConfig := config.NewDBConfig(envConfig)

	productRepository := repository.NewProductRepository()
	categoryRepository := repository.NewCategoryRepository()
	repositoryContainer := NewRepositoryContainer(productRepository, categoryRepository)

	productUseCase := use_case.NewProductUseCase(productDBConfig, productRepository)
	categoryUseCase := use_case.NewCategoryUseCase(productDBConfig, categoryRepository)

	useCaseContainer := NewUseCaseContainer(productUseCase, categoryUseCase)
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productUseCase)
	pb.RegisterCategoryServiceServer(grpcServer, categoryUseCase)

	webContainer := &WebContainer{
		Env:               envConfig,
		ProductDB:         productDBConfig,
		ProductRepository: repositoryContainer,
		UseCase:           useCaseContainer,
		Grpc:              grpcServer,
	}

	return webContainer
}
