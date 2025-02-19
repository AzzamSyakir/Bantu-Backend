package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type ProposalService struct {
	JobRepository *repository.JobRepository
	Producer      *producer.ServicesProducer
}

func NewProposalService(userRepository *repository.JobRepository, producer *producer.ServicesProducer) *ProposalService {
	ProposalService := &ProposalService{
		Producer:      producer,
		JobRepository: userRepository,
	}
	return ProposalService
}
