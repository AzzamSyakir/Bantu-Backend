package controllers

import "bantu-backend/src/internal/services"

type ProposalController struct {
	ProposalService *services.ProposalService
}

func NewProposalController(jobService *services.ProposalService) *ProposalController {
	return &ProposalController{
		ProposalService: jobService,
	}
}
