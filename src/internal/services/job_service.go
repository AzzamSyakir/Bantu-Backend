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
		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseError Failed Get Jobs", err.Error())
	}
	if jobService.JobRepository == nil {
		log.Fatal("JobRepository is nil in GetJobsService")
		return err
	}

	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", getJob)
}

func (jobService *JobService) CreateJobService(request *entity.JobEntity) error {
	createJob, err := jobService.JobRepository.CreateJobRepository(request)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseError ", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) GetJobByIDService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	job, err := jobService.JobRepository.GetJobByIDRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseError Get by id not found", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", job)
}

func (jobService *JobService) UpdateJobService(reader *http.Request, request *entity.JobEntity) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	createJob, err := jobService.JobRepository.UpdateJobRepository(id, request)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseError Update Failed failed", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", createJob)
}

func (jobService *JobService) DeleteJobService(reader *http.Request) error {
	vars := mux.Vars(reader)
	id, _ := vars["id"]
	err := jobService.JobRepository.DeleteJobRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", err.Error())
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", err)
}

// func (jobService *JobService) ApplyJobService(request *entity.ProposalEntity) error {
// 	job, err := jobService.JobRepository.ApplyJobRepository(request)
// 	if err != nil {
// 		return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", err.Error())
// 	}
// 	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "responseSuccess", job)
// }
