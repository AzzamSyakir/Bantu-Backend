package controllers

import (
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"net/http"
	"time"
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
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.LoginService(request)
	select {
	case responseError := <-authController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-authController.ResponseChannel.ResponseSuccess:
		var userEntity entity.UserEntity
		userEntityBytes, err := json.Marshal(responseSuccess.Data)
		if err != nil {
			panic(err.Error())
		}

		err = json.Unmarshal(userEntityBytes, &userEntity)
		if err != nil {
			panic(err.Error())
		}

		authorizationCookie := &http.Cookie{
			Name:     "authorization",
			Value:    userEntity.Token,
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		entityIDCookie := &http.Cookie{
			Name:     "entity_id",
			Value:    userEntity.ID,
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		entityRoleCookie := &http.Cookie{
			Name:     "entity_role",
			Value:    userEntity.Role,
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		http.SetCookie(writer, authorizationCookie)
		http.SetCookie(writer, entityIDCookie)
		http.SetCookie(writer, entityRoleCookie)

		response.NewResponse(writer, &responseSuccess)
	}
}

func (authController *AuthController) AdminRegister(writer http.ResponseWriter, reader *http.Request) {
	request := &request.AdminRegisterRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.AdminRegisterService(request)
	select {
	case responseError := <-authController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-authController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (authController *AuthController) AdminLogin(writer http.ResponseWriter, reader *http.Request) {
	request := &request.AdminLoginRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
		return
	}

	authController.AuthService.AdminLoginService(request)
	select {
	case responseError := <-authController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-authController.ResponseChannel.ResponseSuccess:
		var adminEntity entity.AdminEntity
		adminEntityBytes, err := json.Marshal(responseSuccess.Data)
		if err != nil {
			panic(err.Error())
		}

		err = json.Unmarshal(adminEntityBytes, &adminEntity)
		if err != nil {
			panic(err.Error())
		}

		authorizationCookie := &http.Cookie{
			Name:     "authorization",
			Value:    adminEntity.Token,
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		entityIDCookie := &http.Cookie{
			Name:     "entity_id",
			Value:    adminEntity.ID,
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		entityRoleCookie := &http.Cookie{
			Name:     "entity_role",
			Value:    "admin",
			Expires:  time.Now().Add(24 * time.Hour), // Berlaku 1 hari
			HttpOnly: true,                           // Hanya bisa diakses melalui HTTP, bukan JS
			Secure:   true,                           // Hanya dikirim melalui HTTPS
			Path:     "/",                            // Berlaku di seluruh domain
			Domain:   "localhost",                    // Berlaku di domain localhost
		}

		http.SetCookie(writer, authorizationCookie)
		http.SetCookie(writer, entityIDCookie)
		http.SetCookie(writer, entityRoleCookie)

		response.NewResponse(writer, &responseSuccess)
	}
}
