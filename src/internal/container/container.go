package container

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/middlewares"
	"bantu-backend/src/internal/repository"
	"bantu-backend/src/internal/routes"
	"bantu-backend/src/internal/services"
	"log"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Container struct {
	Env        *configs.EnvConfig
	Db         *configs.DatabaseConfig
	Controller *ControllerContainer
	Route      *routes.Route
	Middleware *middlewares.Middleware
}

func NewContainer() *Container {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	envConfig := configs.NewEnvConfig()
	dbConfig := configs.NewDBConfig(envConfig)
	// setup repo
	userRepository := repository.NewUserRepository()
	chatRepository := repository.NewChatRepository()
	jobRepository := repository.NewJobRepository()
	transactionRepository := repository.NewTransactionRepository()
	// setup services
	authService := services.NewAuthService(userRepository)
	userService := services.NewUserService(userRepository)
	chatService := services.NewChatService(chatRepository)
	jobService := services.NewJobService(jobRepository)
	proposalService := services.NewProposalService(jobRepository)
	transactionService := services.NewTransactionService(transactionRepository)
	// setup controller
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	chatController := controllers.NewChatController(chatService)
	jobController := controllers.NewJobController(jobService)
	proposalController := controllers.NewProposalController(proposalService)
	transactionController := controllers.NewTransactionController(transactionService)
	// setup controllerContainer
	controllerContainer := NewControllerContainer(authController, userController, chatController, jobController, proposalController, transactionController)
	router := mux.NewRouter()
	middleware := middlewares.NewMiddleware()
	routeConfig := routes.NewRoute(
		router,
		middleware,
		authController,
		chatController,
		jobController,
		proposalController,
		transactionController,
	)
	routeConfig.Register()
	container := &Container{
		Db:         dbConfig,
		Controller: controllerContainer,
		Route:      routeConfig,
		Env:        envConfig,
		Middleware: middleware,
	}
	return container
}
