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
	ResponseChannel *response.ResponseChannel
}

func NewAuthController(authService *services.AuthService, responseChannel *response.ResponseChannel) *AuthController {
	return &AuthController{
		AuthService:     authService,
		ResponseChannel: responseChannel,
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
	select {
	case responseError := <-authController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-authController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
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
	select {
	case responseError := <-authController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-authController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}
