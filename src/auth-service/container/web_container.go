package container

import (
	"bantu-backend/grpc/pb"
	"bantu-backend/src/auth-service/config"
	"bantu-backend/src/auth-service/delivery/grpc/client"
	httpdelivery "bantu-backend/src/auth-service/delivery/http"
	"bantu-backend/src/auth-service/delivery/http/middleware"
	"bantu-backend/src/auth-service/delivery/http/route"
	"bantu-backend/src/auth-service/repository"
	"bantu-backend/src/auth-service/use_case"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type WebContainer struct {
	Env        *config.EnvConfig
	AuthDB     *config.DatabaseConfig
	Repository *RepositoryContainer
	UseCase    *UseCaseContainer
	Controller *ControllerContainer
	Route      *route.RootRoute
	Grpc       *grpc.Server
}

func NewWebContainer() *WebContainer {
	errEnvLoad := godotenv.Load()
	if errEnvLoad != nil {
		panic(fmt.Errorf("error loading .env file: %w", errEnvLoad))
	}

	envConfig := config.NewEnvConfig()
	authDBConfig := config.NewAuthDBConfig(envConfig)

	authRepository := repository.NewAuthRepository()
	repositoryContainer := NewRepositoryContainer(authRepository)

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
	orderUrl := fmt.Sprintf(
		"%s:%s",
		envConfig.App.OrderHost,
		envConfig.App.OrderPort,
	)

	initUserClient := client.InitUserServiceClient(userUrl)
	initProductClient := client.InitProductServiceClient(productUrl)

	initOrderClient := client.InitOrderServiceClient(orderUrl)
	initCategoryClient := client.InitCategoryServiceClient(productUrl)
	authUseCase := use_case.NewAuthUseCase(authDBConfig, authRepository, envConfig, &initUserClient)
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authUseCase)

	exposeUseCase := use_case.NewExposeUseCase(authDBConfig, authRepository, envConfig, &initUserClient, &initProductClient, &initOrderClient, &initCategoryClient)

	useCaseContainer := NewUseCaseContainer(authUseCase, exposeUseCase)

	authController := httpdelivery.NewAuthController(authUseCase, exposeUseCase)
	exposeController := httpdelivery.NewExposeController(exposeUseCase)

	controllerContainer := NewControllerContainer(authController, exposeController)

	router := mux.NewRouter()
	authMiddleware := middleware.NewAuthMiddleware(*authRepository, authDBConfig)
	authRoute := route.NewAuthRoute(router, authController)
	// expose route
	userRoute := route.NewUserRoute(router, exposeController, authMiddleware)
	productRoute := route.NewProductRoute(router, exposeController, authMiddleware)
	categoryRoute := route.NewCategoryRoute(router, exposeController, authMiddleware)
	orderRoute := route.NewOrderRoute(router, exposeController, authMiddleware)

	rootRoute := route.NewRootRoute(
		router,
		authRoute,
	)
	exposeRoute := route.NewExposeRoute(
		router,
		userRoute,
		productRoute,
		categoryRoute,
		orderRoute,
	)

	rootRoute.Register()
	exposeRoute.Register()

	webContainer := &WebContainer{
		Env:        envConfig,
		AuthDB:     authDBConfig,
		Repository: repositoryContainer,
		UseCase:    useCaseContainer,
		Controller: controllerContainer,
		Route:      rootRoute,
		Grpc:       grpcServer,
	}

	return webContainer
}
