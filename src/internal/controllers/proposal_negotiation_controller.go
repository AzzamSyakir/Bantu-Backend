package controllers

import (
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProposalController struct {
	ProposalService *services.ProposalService
	ResponseChannel *response.ResponseChannel
}

func NewProposalController(jobService *services.ProposalService, responseChannel *response.ResponseChannel) *ProposalController {
	return &ProposalController{
		ProposalService: jobService,
		ResponseChannel: responseChannel,
	}
}

func (proposalController *ProposalController) GetProposals(writer http.ResponseWriter, reader *http.Request) {

	proposalController.ProposalService.GetProposalsService(reader)
	select {
	case responseError := <-proposalController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-proposalController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (proposalController *ProposalController) CreateProposal(writer http.ResponseWriter, reader *http.Request) {

	request := request.ProposalRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)

	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	proposalController.ProposalService.CreateProposalService(&request)
	select {
	case responseError := <-proposalController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-proposalController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (proposalController *ProposalController) UpdateProposal(writer http.ResponseWriter, reader *http.Request) {

	request := request.ProposalRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)

	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	proposalController.ProposalService.UpdateProposalService(reader, &request)
	select {
	case responseError := <-proposalController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-proposalController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (proposalController *ProposalController) AcceptProposal(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id := vars["proposalId"]
	proposalController.ProposalService.AcceptProposalService(id)
	select {
	case responseError := <-proposalController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-proposalController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}
