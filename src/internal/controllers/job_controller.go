package controllers

import (
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
)

type JobController struct {
	JobService      *services.JobService
	ResponseChannel chan response.Response[any]
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{
		JobService:      jobService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}
