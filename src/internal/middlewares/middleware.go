package middlewares

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/rabbitmq/producer"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Middleware struct {
	Rabbitmq *configs.RabbitMqConfig
	Producer *producer.ServicesProducer
	Cors     *cors.Cors
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Cors: cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
			Debug:            false,
		}),
	}
}
func (middleware *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	limiter := rate.NewLimiter(100, 100)
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

func (m *Middleware) ApplyMiddleware(next http.Handler) http.Handler {
	handler := m.Cors.Handler(next)
	handler = m.RateLimitMiddleware(handler)
	handler = m.InputValidationMiddleware(handler)
	return handler
}
