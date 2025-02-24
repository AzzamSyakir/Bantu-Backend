package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/xendit/xendit-go/v4"
	"github.com/xendit/xendit-go/v4/invoice"
)

type TransactionService struct {
	Env                   *configs.EnvConfig
	TransactionRepository *repository.TransactionRepository
	UserRepository        *repository.UserRepository
	ProposalRepository    *repository.JobRepository
	Producer              *producer.ServicesProducer
	Database              *configs.DatabaseConfig
	Rabbitmq              *configs.RabbitMqConfig
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	userRepository *repository.UserRepository,
	proposalRepository *repository.JobRepository,
	producer *producer.ServicesProducer,
	databaseConfig *configs.DatabaseConfig,
	rabbitmq *configs.RabbitMqConfig,
	env *configs.EnvConfig,
) *TransactionService {
	TransactionService := &TransactionService{
		Env:                   env,
		Database:              databaseConfig,
		Rabbitmq:              rabbitmq,
		Producer:              producer,
		TransactionRepository: transactionRepository,
		UserRepository:        userRepository,
		ProposalRepository:    proposalRepository,
	}
	return TransactionService
}
func (transactionService *TransactionService) TopUpBalance(request *request.TopupRequest, userId string) {
	begin, beginErr := transactionService.Database.DB.Connection.Begin()
	if beginErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, beginErr.Error(), http.StatusInternalServerError)
	}

	if request.Amount == 0 {
		errMessage := "topup amount must be provided"
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}
	newTransaction := &entity.TransactionEntity{
		ID:              string(uuid.NewString()),
		SenderId:        null.NewString(userId, true),
		TransactionType: "top_up",
		Amount:          request.Amount,
		PaymentMethod:   request.PaymentMethod,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if newTransaction.Status == "" {
		newTransaction.Status = "pending"
	}
	_, createdTransactionErr := transactionService.TransactionRepository.CreateTransaction(begin, newTransaction)
	if createdTransactionErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, createdTransactionErr.Error(), http.StatusInternalServerError)
		return
	}
	commitErr := begin.Commit()
	if commitErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}
	// create invoice
	createInvoice := transactionService.CreateInvoiceRequest(userId, float64(newTransaction.Amount))
	newTransaction.InvoiceUrl = createInvoice.InvoiceUrl
	// run function for checking status invoice
	go transactionService.CheckInvoice(userId, createInvoice, newTransaction)
	transactionService.Producer.CreateMessageTransaction(transactionService.Rabbitmq.Channel, newTransaction)
}
func (transactionService *TransactionService) PayFreelancer(request *request.TopupRequest, proposalId string, userId string) {
	begin, beginErr := transactionService.Database.DB.Connection.Begin()
	if beginErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, beginErr.Error(), http.StatusInternalServerError)
	}
	if proposalId == "" {
		errMessage := "proposal id must be provided"
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}
	if request.Amount == 0 {
		errMessage := "topup amount must be provided"
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}
	foundProposal, foundProposalErr := transactionService.ProposalRepository.GetProposalsById(proposalId)
	if foundProposalErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr.Error(), http.StatusInternalServerError)
		return
	}
	foundSender, foundSenderErr := transactionService.UserRepository.GetUserById(begin, userId)
	if foundSenderErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	foundReceiver, foundReceiverErr := transactionService.UserRepository.GetUserById(begin, foundProposal.FreelancerID.String())
	if foundReceiverErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	senderBalance := foundSender.Balance - float64(request.Amount)
	receiverBalance := foundReceiver.Balance + float64(request.Amount)
	// update user balance
	updateSenderBalanceErr := transactionService.UserRepository.UpdateUserBalance(begin, userId, int(senderBalance))
	if updateSenderBalanceErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	updateReceiverBalanceErr := transactionService.UserRepository.UpdateUserBalance(begin, foundProposal.FreelancerID.String(), int(receiverBalance))
	if updateReceiverBalanceErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	newTransaction := &entity.TransactionEntity{
		ID:              string(uuid.NewString()),
		SenderId:        null.NewString(userId, true),
		ReceiverId:      null.NewString(foundProposal.FreelancerID.String(), true),
		TransactionType: "pay_freelancer",
		Amount:          request.Amount,
		PaymentMethod:   request.PaymentMethod,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if newTransaction.Status == "" {
		newTransaction.Status = "pending"
	}
	_, createdTransactionErr := transactionService.TransactionRepository.CreateTransaction(begin, newTransaction)
	if createdTransactionErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, createdTransactionErr.Error(), http.StatusInternalServerError)
		return
	}
	commitErr := begin.Commit()
	if commitErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}
	transactionService.Producer.CreateMessageTransaction(transactionService.Rabbitmq.Channel, newTransaction)
}

func (transactionService *TransactionService) CreateInvoiceRequest(userId string, amount float64) *invoice.Invoice {
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(userId, amount)
	xenditClient := xendit.NewClient(transactionService.Env.XenditSecretKey)
	resp, _, xenditErr := xenditClient.InvoiceApi.CreateInvoice(context.Background()).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()
	if xenditErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, xenditErr.Error(), http.StatusInternalServerError)
		return nil
	}
	return resp
}
func (transactionService *TransactionService) CheckInvoice(userId string, inv *invoice.Invoice, transaction *entity.TransactionEntity) *invoice.Invoice {
	xenditClient := xendit.NewClient(transactionService.Env.XenditSecretKey)
	begin, beginErr := transactionService.Database.DB.Connection.Begin()
	if beginErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, beginErr.Error(), http.StatusInternalServerError)
	}

	for inv.Status != "PAID" && inv.Status != "SETTLED" {
		resp, _, xenditErr := xenditClient.InvoiceApi.GetInvoiceById(context.Background(), *inv.Id).Execute()
		if xenditErr != nil {
			transactionService.Producer.CreateMessageError(
				transactionService.Rabbitmq.Channel,
				xenditErr.Error(),
				http.StatusInternalServerError,
			)
			return nil
		}

		inv.Status = resp.Status
		fmt.Println("status transaction ", inv.Status)
		if inv.Status == "PAID" || inv.Status == "SETTLED" {
			transaction.Status = "completed"
			transactionService.TransactionRepository.UpdateTransactionStatus(begin, transaction)
			transactionService.UserRepository.UpdateUserBalance(begin, userId, transaction.Amount)
			begin.Commit()
			time.Sleep(1 * time.Second)
			break
		}
		time.Sleep(1 * time.Second)
	}

	return inv
}
