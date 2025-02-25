package container

import (
	"bantu-backend/src/cache"
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/controllers"
	"bantu-backend/src/internal/middleware"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/rabbitmq/consumer"
	"bantu-backend/src/internal/rabbitmq/producer"
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
	RabbitMq   *configs.RabbitMqConfig
	Redis      *configs.RedisConfig
	Route      *routes.Route
	Middleware *middleware.Middleware
	WebSocket  *routes.WebSocketServer
}

func NewContainer() *Container {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	envConfig := configs.NewEnvConfig()
	dbConfig := configs.NewDBConfig(envConfig)
	servicesProducer := producer.CreateNewServicesProducer(envConfig.RabbitMq)
	rabbitmqConfig := configs.NewRabbitMqConfig(envConfig)
	redisConfig := configs.NewRedisConfig(envConfig)
	middleware := middleware.NewMiddleware(rabbitmqConfig, servicesProducer, envConfig)

	// setup repo
	userRepository := repository.NewUserRepository()
	chatRepository := repository.NewChatRepository(dbConfig)
	jobRepository := repository.NewJobRepository(dbConfig)
	transactionRepository := repository.NewTransactionRepository()

	// setup cache
	jobCache := cache.NewJobCache(redisConfig)

	// setup services
	authService := services.NewAuthService(userRepository, servicesProducer, envConfig, dbConfig, rabbitmqConfig)
	userService := services.NewUserService(userRepository, servicesProducer)
	chatService := services.NewChatService(chatRepository, servicesProducer, rabbitmqConfig)
	jobService := services.NewJobService(jobRepository, servicesProducer, rabbitmqConfig, jobCache)
	proposalService := services.NewProposalService(jobRepository, servicesProducer, rabbitmqConfig)
	transactionService := services.NewTransactionService(transactionRepository, userRepository, jobRepository, servicesProducer, dbConfig, rabbitmqConfig, envConfig)
	// setup controller
	responseChannel := response.NewResponseChannel()
	authController := controllers.NewAuthController(authService, responseChannel)
	userController := controllers.NewUserController(userService, responseChannel)
	chatController := controllers.NewChatController(chatService, responseChannel)
	jobController := controllers.NewJobController(jobService, responseChannel)
	proposalController := controllers.NewProposalController(proposalService, responseChannel)
	transactionController := controllers.NewTransactionController(transactionService, responseChannel)
	// setup controllerContainer
	controllerContainer := NewControllerContainer(authController, userController, chatController, jobController, proposalController, transactionController)
	controllerConsumer := consumer.NewControllerConsumer(envConfig.RabbitMq, authController, chatController, jobController, proposalController, transactionController, userController, responseChannel)
	consumerInit := consumer.NewConsumerEntrypointInit(controllerConsumer, rabbitmqConfig)
	consumerInit.ConsumerEntrypointStart()
	router := mux.NewRouter()
	websocket := routes.NewWebSocketServer(router, chatController)
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
		Env:        envConfig,
		Db:         dbConfig,
		Controller: controllerContainer,
		RabbitMq:   rabbitmqConfig,
		Redis:      redisConfig,
		Route:      routeConfig,
		Middleware: middleware,
		WebSocket:  websocket,
	}
	return container
}
