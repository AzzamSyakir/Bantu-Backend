package services

import (
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
)

type JobService struct {
	JobRepository *repository.JobRepository
	Producer      *producer.ServicesProducer
}

func NewJobService(jobRepo *repository.JobRepository, producer *producer.ServicesProducer) *JobService {
	return &JobService{
		Producer:      producer,
		JobRepository: jobRepo,
	}
}
