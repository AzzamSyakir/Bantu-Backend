package controllers

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/models/response"
	"bantu-backend/src/internal/services"
	"encoding/json"
	"net/http"
)

type JobController struct {
	JobService      *services.JobService
	Rabbitmq        *configs.RabbitMqConfig
	ResponseChannel *response.ResponseChannel
}

func NewJobController(jobService *services.JobService, responseChannel *response.ResponseChannel) *JobController {
	return &JobController{
		JobService:      jobService,
		ResponseChannel: responseChannel,
	}
}

func (jobController *JobController) GetJobs(writer http.ResponseWriter, reader *http.Request) {
	jobController.JobService.GetJobsService(writer, reader)
	select {
	case responseError := <-jobController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-jobController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (jobController *JobController) CreateJob(writer http.ResponseWriter, reader *http.Request) {
	request := &request.JobRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	jobController.JobService.CreateJobService(request)
	select {
	case responseError := <-jobController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-jobController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (jobController *JobController) GetJobByID(writer http.ResponseWriter, reader *http.Request) {

	jobController.JobService.GetJobByIDService(reader)
	select {
	case responseError := <-jobController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-jobController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (jobController *JobController) UpdateJob(writer http.ResponseWriter, reader *http.Request) {

	request := &request.JobRequest{}
	decodeErr := json.NewDecoder(reader.Body).Decode(&request)
	if decodeErr != nil {
		http.Error(writer, decodeErr.Error(), 404)
	}

	jobController.JobService.UpdateJobService(reader, request)
	select {
	case responseError := <-jobController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-jobController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}

func (jobController *JobController) DeleteJob(writer http.ResponseWriter, reader *http.Request) {

	jobController.JobService.DeleteJobService(reader)
	select {
	case responseError := <-jobController.ResponseChannel.ResponseError:
		response.NewResponse(writer, &responseError)
	case responseSuccess := <-jobController.ResponseChannel.ResponseSuccess:
		response.NewResponse(writer, &responseSuccess)
	}
}
