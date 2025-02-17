package controllers

import "bantu-backend/src/internal/services"

type JobController struct {
	JobService *services.JobService
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{
		JobService: jobService,
	}
}
