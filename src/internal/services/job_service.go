package services

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
)

type JobService struct {
	JobRepository *repository.JobRepository
	RabbitMq      *configs.RabbitMqConfig
	Producer      *producer.ServicesProducer
}

func NewJobService(jobRepo *repository.JobRepository, producer *producer.ServicesProducer, rabbitmq *configs.RabbitMqConfig) *JobService {
	return &JobService{
		JobRepository: jobRepo,
		Producer:      producer,
		RabbitMq:      rabbitmq,
	}
}

func (jobService *JobService) GetJobsService(writer http.ResponseWriter, reader *http.Request) error {
	queryParams := reader.URL.Query()
	getJob, err := jobService.JobRepository.GetJobsRepository(queryParams)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}

	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", getJob)
}

func (jobService *JobService) CreateJobService(request *request.JobRequest) error {
	job := &entity.JobEntity{
		Title:       request.Title,
		Description: request.Description,
		Category:    request.Category,
		Price:       request.Price,
		RegencyID:   request.RegencyID,
		ProvinceID:  request.ProvinceID,
	}
	createJob, err := jobService.JobRepository.CreateJobRepository(job)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "create job is failed", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) GetJobByIDService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	job, err := jobService.JobRepository.GetJobByIDRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "job not found", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", job)
}

func (jobService *JobService) UpdateJobService(reader *http.Request, request *request.JobRequest) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	job := &entity.JobEntity{
		Title:       request.Title,
		Description: request.Description,
		Category:    request.Category,
		Price:       request.Price,
		RegencyID:   request.RegencyID,
		ProvinceID:  request.ProvinceID,
	}
	createJob, err := jobService.JobRepository.UpdateJobRepository(id, job)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "update job is failed", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) DeleteJobService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	err := jobService.JobRepository.DeleteJobRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "delete job is failed", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", "Success delete job")
}
