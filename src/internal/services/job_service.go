package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type JobService struct {
	Database      *configs.DatabaseConfig
	JobRepository *repository.JobRepository
	Rabbitmq      *configs.RabbitMqConfig
	Producer      *producer.ServicesProducer
}

func NewJobService(jobRepo *repository.JobRepository, producer *producer.ServicesProducer, database *configs.DatabaseConfig, rabbitmq *configs.RabbitMqConfig) *JobService {
	return &JobService{
		Database:      database,
		JobRepository: jobRepo,
		Producer:      producer,
		Rabbitmq:      rabbitmq,
	}
}

func (jobService *JobService) GetJobsService(writer http.ResponseWriter, reader *http.Request) error {
	getJob, err := jobService.JobRepository.GetJobsRepository()
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseError Failed Get Jobs", err.Error())
	}
	if jobService.JobRepository == nil {
		log.Fatal("JobRepository is nil in GetJobsService")
	}

	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", getJob)
}

func (jobService *JobService) CreateJobService(request *entity.JobEntity) error {
	createJob, err := jobService.JobRepository.CreateJobRepository(request)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseError ", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) GetJobByIDService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	job, err := jobService.JobRepository.GetJobByIDRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseError Database connection failed", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", job)
}

func (jobService *JobService) UpdateJobService(reader *http.Request, request *entity.JobEntity) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	createJob, err := jobService.JobRepository.UpdateJobRepository(id, request)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseError Database connection failed", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) DeleteJobService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	err := jobService.JobRepository.DeleteJobRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", err)
}

func (jobService *JobService) ApplyJobService(request *entity.ProposalEntity) error {
	job, err := jobService.JobRepository.ApplyJobRepository(request)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.Rabbitmq.Channel, "responseSuccess", job)
}
