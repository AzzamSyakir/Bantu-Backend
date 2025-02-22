package middlewares

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/services"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Middleware struct {
	Rabbitmq *configs.RabbitMqConfig
	Producer *producer.ServicesProducer
	Cors     *cors.Cors
	Auth     *services.AuthService
}

func NewMiddleware(auth *services.AuthService) *Middleware {
	return &Middleware{
		Cors: cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
			Debug:            false,
		}),
		Auth: auth,
	}
}
func (middleware *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(2, 4)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {

			errorMessage := "The API is at capacity, try again later."
			w.WriteHeader(http.StatusTooManyRequests)
			middleware.Producer.CreateMessageError(middleware.Rabbitmq.Channel, errorMessage, http.StatusTooManyRequests)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (m *Middleware) InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ValidateAuthorizationHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if m.Auth.ValidateToken(r.Header.Get("Authorization")) == false {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) ApplyMiddleware(next http.Handler) http.Handler {
	handler := m.Cors.Handler(next)
	handler = m.RateLimitMiddleware(handler)
	handler = m.InputValidationMiddleware(handler)
	handler = m.ValidateAuthorizationHeader(handler)
	return handler
}
