package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProposalService struct {
	JobRepository *repository.JobRepository
	RabbitMq      *configs.RabbitMqConfig
	Producer      *producer.ServicesProducer
}

func NewProposalService(userRepository *repository.JobRepository, producer *producer.ServicesProducer, rabbitMq *configs.RabbitMqConfig) *ProposalService {
	ProposalService := &ProposalService{
		Producer:      producer,
		JobRepository: userRepository,
		RabbitMq:      rabbitMq,
	}
	return ProposalService
}

func (proposalService *ProposalService) GetProposalsService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	proposals, err := proposalService.JobRepository.GetProposalsRepository(id)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "get proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, proposals)
}

func (proposalService *ProposalService) CreateProposalService(request *request.ProposalRequest) error {
	proposal := &entity.ProposalEntity{
		ID:            uuid.NewString(),
		JobID:         request.JobID,
		FreelancerID:  request.FreelancerID,
		ProposalText:  request.ProposalText,
		ProposedPrice: request.ProposedPrice,
		Status:        request.Status,
	}
	proposal, err := proposalService.JobRepository.CreateProposalRepository(proposal)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "create proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, proposal)
}

func (proposalService *ProposalService) UpdateProposalService(reader *http.Request, request *request.ProposalRequest) error {
	vars := mux.Vars(reader)
	id := vars["proposalId"]
	proposal := &entity.ProposalEntity{
		JobID:         request.JobID,
		FreelancerID:  request.FreelancerID,
		ProposalText:  request.ProposalText,
		ProposedPrice: request.ProposedPrice,
	}
	proposal, err := proposalService.JobRepository.UpdateProposalRepository(id, proposal)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "update proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, proposal)
}

func (proposalService *ProposalService) AcceptProposalService(id string) error {
	_, err := proposalService.JobRepository.AcceptProposalRepository(id)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageJob(proposalService.RabbitMq.Channel, "success accept proposal")
}
