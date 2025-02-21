package controllers

import (
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type ProposalController struct {
	ProposalService *services.ProposalService
	ResponseChannel chan response.Response[any]
}

func NewProposalController(jobService *services.ProposalService) *ProposalController {
	return &ProposalController{
		ProposalService: jobService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}

func (proposalController *ProposalController) GetProposals(writer http.ResponseWriter, reader *http.Request) {

	proposalController.ProposalService.GetProposalsService(reader)
	responseData := <-proposalController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (proposalController *ProposalController) CreateProposal(writer http.ResponseWriter, reader *http.Request) {

	request := entity.ProposalEntity{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)

	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	proposalController.ProposalService.CreateProposalService(&request)
	responseData := <-proposalController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (proposalController *ProposalController) UpdateProposal(writer http.ResponseWriter, reader *http.Request) {

	request := entity.ProposalEntity{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)

	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	proposalController.ProposalService.UpdateProposalService(reader, &request)
	responseData := <-proposalController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (proposalController *ProposalController) AcceptProposal(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	id, _ := vars["proposalId"]
	proposalController.ProposalService.AcceptProposalService(id)
	responseData := <-proposalController.ResponseChannel
	response.NewResponse(writer, &responseData)
}
