package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	xendit "github.com/xendit/xendit-go/v6"
	invoice "github.com/xendit/xendit-go/v6/invoice"
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
	createInvoice := transactionService.CreateInvoiceRequest(userId, float64(newTransaction.Amount))
	newTransaction.InvoiceUrl = null.NewString(createInvoice.InvoiceUrl, true)
	go transactionService.CheckInvoice(userId, createInvoice, newTransaction)
	transactionService.Producer.CreateMessageTransaction(transactionService.Rabbitmq.Channel, newTransaction)
}
func (transactionService *TransactionService) WithdrawBalance(request *request.WithdrawRequest, userId string) {
	begin, beginErr := transactionService.Database.DB.Connection.Begin()
	if beginErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, beginErr.Error(), http.StatusInternalServerError)
		return
	}
	if request.Amount == 0 {
		errMessage := "withdraw amount must be provided"
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, errMessage, http.StatusInternalServerError)
		return
	}
	userFound, foundUserErr := transactionService.UserRepository.GetUserById(begin, userId)
	if foundUserErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundUserErr.Error(), http.StatusInternalServerError)
		return
	}
	newTransaction := &entity.TransactionEntity{
		ID:              uuid.NewString(),
		ReceiverId:      null.NewString(userId, true),
		TransactionType: "withdrawal",
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
	finalUserBalance := userFound.Balance - float64(newTransaction.Amount)
	commitErr := begin.Commit()
	if commitErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, commitErr.Error(), http.StatusInternalServerError)
		return
	}
	// create payout request
	createPayout := transactionService.CreatePayoutRequest(newTransaction, request)
	go transactionService.CheckPayout(newTransaction, createPayout.ID, int(finalUserBalance))
	transactionService.Producer.CreateMessageTransaction(transactionService.Rabbitmq.Channel, createPayout)
}
func (transactionService *TransactionService) PayFreelancer(request *request.PayFreelancerRequest, proposalId string, userId string) {
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
	foundReceiver, foundReceiverErr := transactionService.UserRepository.GetUserById(begin, foundProposal.FreelancerID)
	if foundReceiverErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	if foundSender.Balance < *foundProposal.ProposedPrice {
		begin.Rollback()
		errMsg := fmt.Sprintf("Insufficient funds: your current balance is %.2f, but the required amount is %.2f.", foundSender.Balance, *foundProposal.ProposedPrice)
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, errMsg, http.StatusInternalServerError)
		return
	}
	senderBalance := foundSender.Balance - float64(*foundProposal.ProposedPrice)
	receiverBalance := foundReceiver.Balance + float64(*foundProposal.ProposedPrice)
	updateSenderBalanceErr := transactionService.UserRepository.UpdateUserBalance(begin, userId, int(senderBalance))
	if updateSenderBalanceErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	updateReceiverBalanceErr := transactionService.UserRepository.UpdateUserBalance(begin, foundProposal.FreelancerID, int(receiverBalance))
	if updateReceiverBalanceErr != nil {
		begin.Rollback()
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, foundProposalErr, http.StatusInternalServerError)
		return
	}
	newTransaction := &entity.TransactionEntity{
		ID:              string(uuid.NewString()),
		SenderId:        null.NewString(userId, true),
		ReceiverId:      null.NewString(foundProposal.FreelancerID, true),
		TransactionType: "pay_freelancer",
		Amount:          int(*foundProposal.ProposedPrice),
		PaymentMethod:   request.PaymentMethod,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if newTransaction.Status == "" {
		newTransaction.Status = "completed"
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
func (transactionService *TransactionService) CreatePayoutRequest(transaction *entity.TransactionEntity, request *request.WithdrawRequest) *entity.TransactionEntity {
	endpoint := "https://api.xendit.co/payouts"
	payload := map[string]interface{}{
		"external_id": transaction.ID,
		"amount":      request.Amount,
		"email":       request.Email,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, fmt.Sprintf("Failed to marshal JSON: %v", err), http.StatusInternalServerError)
		return &entity.TransactionEntity{}
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return &entity.TransactionEntity{}
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(transactionService.Env.XenditSecretKey, "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return &entity.TransactionEntity{}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return &entity.TransactionEntity{}
	}
	var transactionResponse *entity.TransactionEntity
	err = json.Unmarshal(body, &transactionResponse)
	if err != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, err.Error(), http.StatusInternalServerError)
		return &entity.TransactionEntity{}
	}

	return transactionResponse
}

func (transactionService *TransactionService) CheckInvoice(userId string, inv *invoice.Invoice, transaction *entity.TransactionEntity) {
	xenditClient := xendit.NewClient(transactionService.Env.XenditSecretKey)
	begin, beginErr := transactionService.Database.DB.Connection.Begin()
	if beginErr != nil {
		transactionService.Producer.CreateMessageError(transactionService.Rabbitmq.Channel, beginErr.Error(), http.StatusInternalServerError)
	}

	for inv.Status != "PAID" && inv.Status != "SETTLED" {
		resp, _, xenditErr := xenditClient.InvoiceApi.GetInvoiceById(context.Background(), *inv.Id).Execute()
		if xenditErr != nil {
			log.Fatal(xenditErr.Error())
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

}
func (transactionService *TransactionService) CheckPayout(transaction *entity.TransactionEntity, payoutId string, finalUserBalance int) {
	endpoint := "https://api.xendit.co/payouts/" + payoutId
	client := &http.Client{}

	begin, err := transactionService.Database.DB.Connection.Begin()
	if err != nil {
		panic(err)
	}

	for transaction.Status != "PAID" && transaction.Status != "SETTLED" {
		fmt.Println("status transaction", transaction.Status)
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			panic(err)
		}
		req.SetBasicAuth(transactionService.Env.XenditSecretKey, "")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			panic(err)
		}

		var payoutResponse struct {
			Status string `json:"status"`
		}
		err = json.Unmarshal(body, &payoutResponse)
		if err != nil {
			panic(err)
		}

		transaction.Status = payoutResponse.Status

		if transaction.Status == "PAID" || transaction.Status == "SETTLED" || transaction.Status == "COMPLETED" {
			transaction.Status = "completed"
			if err := transactionService.TransactionRepository.UpdateTransactionStatus(begin, transaction); err != nil {
				begin.Rollback()
				panic(err)
			}
			if err := transactionService.UserRepository.UpdateUserBalance(begin, transaction.ReceiverId.String, finalUserBalance); err != nil {
				begin.Rollback()
				panic(err)
			}
			begin.Commit()
			time.Sleep(1 * time.Second)
			break
		}
		time.Sleep(1 * time.Second)
	}
}
