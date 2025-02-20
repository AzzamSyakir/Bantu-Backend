package controllers

import (
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type TransactionController struct {
	TransactionService *services.TransactionService
	ResponseChannel    *response.ResponseChannel
}

func NewTransactionController(jobService *services.TransactionService, responseChannel *response.ResponseChannel) *TransactionController {
	return &TransactionController{
		TransactionService: jobService,
		ResponseChannel:    responseChannel,
	}
}
func (transactionController *TransactionController) TopUpBalance(writer http.ResponseWriter, reader *http.Request) {
	request := &request.TopupRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		responseError := response.Response[any]{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
			Data:    decodeErr.Error(),
		}
		response.NewResponse(writer, &responseError)
		return
	}
	cookie, err := reader.Cookie("entity_id")
	if err != nil {
		errorMessage := fmt.Sprintln("error getting cookies : ", err.Error())
		responseError := response.Response[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "Unauthorized",
			Data:    errorMessage,
		}
		response.NewResponse(writer, &responseError)
	}
	userId := cookie.Value
	transactionController.TransactionService.TopUpBalance(request, userId)
	select {
	case responseError := <-transactionController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-transactionController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (transactionController *TransactionController) PayFreelancer(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	proposalId := vars["proposalId"]
	request := &request.PayFreelancerRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		responseError := response.Response[any]{
			Code:    http.StatusBadRequest,
			Message: "invalid request body",
			Data:    decodeErr.Error(),
		}
		response.NewResponse(writer, &responseError)
		return
	}
	cookie, err := reader.Cookie("entity_id")
	if err != nil {
		errorMessage := fmt.Sprintln("error getting cookies : ", err.Error())
		responseError := response.Response[interface{}]{
			Code:    http.StatusInternalServerError,
			Message: "Unauthorized",
			Data:    errorMessage,
		}
		response.NewResponse(writer, &responseError)
	}
	userId := cookie.Value
	transactionController.TransactionService.PayFreelancer(request, proposalId, userId)
	select {
	case responseError := <-transactionController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-transactionController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}
