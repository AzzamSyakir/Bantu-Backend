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
	"github.com/xendit/xendit-go/v4"
	"github.com/xendit/xendit-go/v4/invoice"
)

type TransactionService struct {
	Env                   *configs.EnvConfig
	TransactionRepository *repository.TransactionRepository
	UserRepository        *repository.UserRepository
	Producer              *producer.ServicesProducer
	Database              *configs.DatabaseConfig
	Rabbitmq              *configs.RabbitMqConfig
}

func NewTransactionService(
	transactionRepository *repository.TransactionRepository,
	userRepository *repository.UserRepository,
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
		UserId:          userId,
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
