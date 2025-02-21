package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"

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
	job, err := proposalService.JobRepository.GetProposalsRepository(id)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "get proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, "responseSuccess", job)
}

func (proposalService *ProposalService) CreateProposalService(request *entity.ProposalEntity) error {
	proposal, err := proposalService.JobRepository.CreateProposalRepository(request)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "create proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, "responseSuccess", proposal)
}

func (proposalService *ProposalService) UpdateProposalService(reader *http.Request, request *entity.ProposalEntity) error {
	vars := mux.Vars(reader)
	id, _ := vars["proposalId"]
	proposal, err := proposalService.JobRepository.UpdateProposalRepository(id, request)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "update proposal is failed", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageProposal(proposalService.RabbitMq.Channel, "responseSuccess", proposal)
}

func (proposalService *ProposalService) AcceptProposalService(id string) error {
	_, err := proposalService.JobRepository.AcceptProposalRepository(id)
	if err != nil {
		return proposalService.Producer.CreateMessageError(proposalService.RabbitMq.Channel, "accept proposal not found", http.StatusBadRequest)
	}
	return proposalService.Producer.CreateMessageJob(proposalService.RabbitMq.Channel, "responseSuccess", "success accept proposal")
}
