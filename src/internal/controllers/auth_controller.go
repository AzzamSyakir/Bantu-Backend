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

	authController.AuthService.RegisterService(request)
	responseData := <-authController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (authController *AuthController) Login(writer http.ResponseWriter, reader *http.Request) {
	request := &request.LoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		log.Println(decodeErr)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.LoginService(request)
	responseData := <-authController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (authController *AuthController) AdminRegister(writer http.ResponseWriter, reader *http.Request) {
	request := &request.AdminRegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		log.Println(decodeErr)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.AdminRegisterService(request)
	responseData := <-authController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (authController *AuthController) AdminLogin(writer http.ResponseWriter, reader *http.Request) {
	request := &request.AdminLoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		log.Println(decodeErr)
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.AdminLoginService(request)
	responseData := <-authController.ResponseChannel
	response.NewResponse(writer, &responseData)
}
