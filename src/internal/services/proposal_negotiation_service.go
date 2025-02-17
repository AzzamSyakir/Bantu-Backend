package services

import "bantu-backend/src/internal/repository"

type ProposalService struct {
	JobRepository *repository.JobRepository
}

func NewProposalService(userRepository *repository.JobRepository) *ProposalService {
	ProposalService := &ProposalService{
		JobRepository: userRepository,
	}
	return ProposalService
}
