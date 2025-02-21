package controllers

import (
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

type AuthController struct {
	AuthService     *services.AuthService
	ResponseChannel chan response.Response[any]
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService:     authService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}

func (authController *AuthController) Register(writer http.ResponseWriter, reader *http.Request) {
	request := &request.RegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		log.Println(decodeErr)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	service, err := authController.AuthService.RegisterService(request)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.NewResponse(writer, &response.Response[any]{
		Code:    http.StatusOK,
		Message: "Register success",
		Data:    service,
	})
}

func (authController *AuthController) Login(writer http.ResponseWriter, reader *http.Request) {
	request := &request.LoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		log.Println(decodeErr)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	service, err := authController.AuthService.LoginService(request)
	if err != nil {
		log.Println(err)
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response.NewResponse(writer, &response.Response[any]{
		Code:    http.StatusOK,
		Message: "Login success",
		Data:    service,
	})
}
