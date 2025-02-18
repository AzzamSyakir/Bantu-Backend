package controllers

import (
	"bantu-backend/src/internal/helper/response"
	"bantu-backend/src/internal/services"
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
