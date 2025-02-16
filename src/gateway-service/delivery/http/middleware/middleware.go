package middleware

import (
	"bantu-backend/src/gateway-service/config"
	"bantu-backend/src/gateway-service/model/response"
	"bantu-backend/src/gateway-service/repository"
	"net/http"
	"strings"
	"time"

	"github.com/guregu/null"
)

type GatewayMiddleware struct {
	SessionRepository *repository.GatewayRepository
	DatabaseConfig    *config.DatabaseConfig
}

func NewGatewayMiddleware(sessionRepository repository.GatewayRepository, databaseConfig *config.DatabaseConfig) *GatewayMiddleware {
	return &GatewayMiddleware{
		SessionRepository: &sessionRepository,
		DatabaseConfig:    databaseConfig,
	}
}

func (authMiddleware *GatewayMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Gatewayorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Missing token",
			}
			response.NewResponse(w, result)
			return
		}

		begin, err := authMiddleware.DatabaseConfig.GatewayDB.Connection.Begin()
		if err != nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusInternalServerError,
				Message: "transaction error",
			}
			response.NewResponse(w, result)
			return
		}

		session, err := authMiddleware.SessionRepository.FindOneByAccToken(begin, token)
		if err != nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: token not found",
			}
			response.NewResponse(w, result)
			return
		}
		if session == nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Invalid Token",
			}
			response.NewResponse(w, result)
			return
		}
		if session.AccessTokenExpiredAt == null.NewTime(time.Now(), true) {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Token expired",
			}
			response.NewResponse(w, result)
			return
		}
		begin.Commit()
		next.ServeHTTP(w, r)
	})
}
