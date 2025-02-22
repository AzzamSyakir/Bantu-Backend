package middleware

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/rabbitmq/producer"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Middleware struct {
	Rabbitmq *configs.RabbitMqConfig
	Producer *producer.ServicesProducer
}

func NewMiddleware() *Middleware {
	return &Middleware{}
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			errorMessage := "The API is at capacity, try again later."
			http.Error(w, errorMessage, http.StatusTooManyRequests)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (*Middleware) InputValidationMiddleware(next http.Handler) http.Handler {
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

func (middleware *Middleware) ApplyMiddleware(next http.Handler) http.Handler {
	handler := middleware.CorsMiddleware(next)
	handler = middleware.RateLimitMiddleware(handler)
	handler = middleware.InputValidationMiddleware(handler)
	fmt.Println("middleware applied")
	return handler
}
