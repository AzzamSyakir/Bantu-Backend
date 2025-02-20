package middleware

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/rabbitmq/producer"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Middleware struct {
	Rabbitmq  *configs.RabbitMqConfig
	Producer  *producer.ServicesProducer
	EnvConfig *configs.EnvConfig
}

func NewMiddleware(rabbimtMq *configs.RabbitMqConfig, producer *producer.ServicesProducer, envConfig *configs.EnvConfig) *Middleware {
	return &Middleware{
		Rabbitmq:  rabbimtMq,
		Producer:  producer,
		EnvConfig: envConfig,
	}
}

func (*Middleware) CorsMiddleware(next http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            false,
	})
	return c.Handler(next)
}

func (middleware *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(100, 100)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			errorMessage := "The API is at capacity, try again later."
			result := &response.Response[interface{}]{
				Code:    http.StatusTooManyRequests,
				Message: "Error",
				Data:    errorMessage,
			}
			response.NewResponse(writer, result)
			return
		} else {
			next.ServeHTTP(writer, request)
		}
	})
}

func (middleware *Middleware) InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost || request.Method == http.MethodPut {
			if !strings.Contains(request.Header.Get("Content-Type"), "application/json") {
				errorMessage := "Content-Type must be application/json"
				result := &response.Response[interface{}]{
					Code:    http.StatusUnsupportedMediaType,
					Message: "Error",
					Data:    errorMessage,
				}
				response.NewResponse(writer, result)
				return
			}
		}
		next.ServeHTTP(writer, request)
	})
}

func (m *Middleware) ValidateAuthorizationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "login") || strings.Contains(r.URL.Path, "register") {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if !m.ValidateToken(r.Header.Get("Authorization")) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ValidateToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(m.EnvConfig.SecretKey)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return false
	}

	if !token.Valid {
		return false
	}

	return true
}

func (middleware *Middleware) ApplyMiddleware(next http.Handler) http.Handler {
	handler := middleware.CorsMiddleware(next)
	handler = middleware.RateLimitMiddleware(handler)
	handler = middleware.InputValidationMiddleware(handler)
	handler = middleware.ValidateAuthorizationHeader(handler)
	fmt.Println("middleware applied")
	return handler
}
