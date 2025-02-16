package http

import (
	model_request "bantu-backend/src/gateway-service/model/request/controller"
	"bantu-backend/src/gateway-service/model/response"
	"bantu-backend/src/gateway-service/use_case"
	"encoding/json"
	"net/http"
	"strings"
)

type GatewayController struct {
	GatewayUseCase *use_case.GatewayUseCase
	ExposeUseCase  *use_case.ExposeUseCase
}

func NewGatewayController(authUseCase *use_case.GatewayUseCase, exposeUseCase *use_case.ExposeUseCase) *GatewayController {
	authController := &GatewayController{
		GatewayUseCase: authUseCase,
		ExposeUseCase:  exposeUseCase,
	}
	return authController
}
func (authController *GatewayController) Register(writer http.ResponseWriter, reader *http.Request) {

	request := &model_request.RegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	result := authController.ExposeUseCase.CreateUser(request)

	response.NewResponse(writer, result)
}
func (authController *GatewayController) Login(writer http.ResponseWriter, reader *http.Request) {
	request := &model_request.LoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}
	foundUser, _ := authController.GatewayUseCase.Login(request)
	response.NewResponse(writer, foundUser)
}
func (authController *GatewayController) Logout(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Gatewayorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)

	result, _ := authController.GatewayUseCase.Logout(tokenString)
	response.NewResponse(writer, result)
}

func (authController *GatewayController) GetNewAccessToken(writer http.ResponseWriter, reader *http.Request) {
	token := reader.Header.Get("Gatewayorization")
	tokenString := strings.Replace(token, "Bearer ", "", 1)

	result, _ := authController.GatewayUseCase.GetNewAccessToken(tokenString)
	response.NewResponse(writer, result)
}
