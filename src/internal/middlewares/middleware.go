package middlewares

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Middleware struct {
	Cors    *cors.Cors
	limiter *rate.Limiter
}

func NewMiddleware() *Middleware {
	return &Middleware{
		Cors: cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
			AllowedHeaders:   []string{"*"},
			Debug:            false,
		}),
		limiter: rate.NewLimiter(10, 20),
	}
}

func (m *Middleware) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
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
