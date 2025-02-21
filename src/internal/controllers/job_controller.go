package controllers

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"net/http"
)

type JobController struct {
	JobService      *services.JobService
	Rabbitmq        *configs.RabbitMqConfig
	Producer        *producer.ControllerProducer
	ResponseChannel chan response.Response[any]
}

func NewJobController(jobService *services.JobService) *JobController {
	return &JobController{
		JobService:      jobService,
		ResponseChannel: make(chan response.Response[any], 1),
	}
}

func (jobController *JobController) GetJobs(writer http.ResponseWriter, reader *http.Request) {
	jobController.JobService.GetJobsService(writer, reader)
	responseData := <-jobController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (jobController *JobController) CreateJob(writer http.ResponseWriter, reader *http.Request) {
	request := &request.JobRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	jobController.JobService.CreateJobService(request)
	responseData := <-jobController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (jobController *JobController) GetJobByID(writer http.ResponseWriter, reader *http.Request) {

	jobController.JobService.GetJobByIDService(reader)
	responseData := <-jobController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (jobController *JobController) UpdateJob(writer http.ResponseWriter, reader *http.Request) {

	request := &request.JobRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	jobController.JobService.UpdateJobService(reader, request)
	responseData := <-jobController.ResponseChannel
	response.NewResponse(writer, &responseData)
}

func (jobController *JobController) DeleteJob(writer http.ResponseWriter, reader *http.Request) {

	jobController.JobService.DeleteJobService(reader)
	responseData := <-jobController.ResponseChannel
	response.NewResponse(writer, &responseData)
}
