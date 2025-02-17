package services

import "bantu-backend/src/internal/repository"

type JobService struct {
	JobRepository *repository.JobRepository
}

func NewJobService(jobRepo *repository.JobRepository) *JobService {
	return &JobService{
		JobRepository: jobRepo,
	}
}
