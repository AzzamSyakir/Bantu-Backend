package middleware

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/rabbitmq/producer"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		if strings.Contains(reader.URL.Path, "login") || strings.Contains(reader.URL.Path, "register") {
			next.ServeHTTP(writer, reader)
			return
		}

		var tokenString string
		cookie, err := reader.Cookie("authorization")
		if err == nil {
			tokenString = cookie.Value
		} else {
			tokenString = reader.Header.Get("Authorization")
		}

		if tokenString == "" {
			errorMessage := "Unauthorized access: please log in to obtain valid credentials."
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data:    errorMessage,
			}
			response.NewResponse(writer, result)
			return
		}

		log.Println(tokenString)
		tokenDecoded, ValidateTokenOk := m.ValidateToken(tokenString)
		if !ValidateTokenOk {
			errorMessage := "Invalid token: please provide a valid token or log in again."
			result := &response.Response[any]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data:    errorMessage,
			}
			response.NewResponse(writer, result)
			return
		}

		log.Println(reader.URL.Path, reader.Method, tokenDecoded.Rl)
		ValidateRoleOk := m.ValidateRole(reader.URL.Path, reader.Method, tokenDecoded.Rl)
		if !ValidateRoleOk {
			errorMessage := "Unauthorized access: you do not have the required permissions to access this resource."
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized",
				Data:    errorMessage,
			}
			response.NewResponse(writer, result)
			return
		}

		log.Println("Done ValidateAuthorizationHeader")
		next.ServeHTTP(writer, reader)
	})
}

func (m *Middleware) ValidateToken(tokenString string) (*request.Authorization, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		secretKey := []byte(m.EnvConfig.SecretKey)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, false
	}

	id, idOk := claims["id"].(string)
	rl, rlOk := claims["rl"].(string)
	exp, expOk := claims["exp"].(float64)

	if !idOk || !rlOk || !expOk {
		return nil, false
	}

	dataClaims := &request.Authorization{
		Id:  id,
		Rl:  rl,
		Exp: exp,
	}

	return dataClaims, true
}

func (m *Middleware) ValidateRole(endpoint string, method string, role string) bool {
	anyId := `[a-zA-Z0-9_-]+`

	rolePermissions := map[string]map[string][]string{
		`^/api/jobs$`: {
			"GET":  {"freelancer", "company", "client"},
			"POST": {"company", "client"},
		},
		fmt.Sprintf(`^/api/jobs/%s$`, anyId): {
			"GET":    {"freelancer"},
			"PUT":    {"company", "client"},
			"DELETE": {"company", "client"},
		},
		fmt.Sprintf(`^/api/jobs/%s/proposals$`, anyId): {
			"GET": {"company", "client"},
		},
		fmt.Sprintf(`^/api/jobs/%s/proposal$`, anyId): {
			"POST": {"freelancer"},
		},
		fmt.Sprintf(`^/api/jobs/%s/proposal/%s$`, anyId, anyId): {
			"PUT": {"freelancer"},
		},
		fmt.Sprintf(`^/api/jobs/%s/proposal/%s/accept$`, anyId, anyId): {
			"PUT": {"company", "client"},
		},
	}

	for pattern, methodRoles := range rolePermissions {
		matched, _ := regexp.MatchString(pattern, endpoint)
		if matched {
			allowedRoles, ok := methodRoles[method]
			if !ok {
				return false
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					return true
				}
			}
			return false
		}
	}

	return false
}

func (middleware *Middleware) ApplyMiddleware(next http.Handler) http.Handler {
	handler := middleware.CorsMiddleware(next)
	handler = middleware.RateLimitMiddleware(handler)
	handler = middleware.InputValidationMiddleware(handler)
	handler = middleware.ValidateAuthorizationHeader(handler)
	fmt.Println("middleware applied")
	return handler
}
