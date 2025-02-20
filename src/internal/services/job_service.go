package services

import (
	"bantu-backend/src/cache"
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"bantu-backend/src/internal/models/request"
	"bantu-backend/src/internal/rabbitmq/producer"
	"bantu-backend/src/internal/repository"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type JobService struct {
	JobRepository *repository.JobRepository
	JobCache      *cache.JobCache
	RabbitMq      *configs.RabbitMqConfig
	Producer      *producer.ServicesProducer
}

func NewJobService(jobRepo *repository.JobRepository, producer *producer.ServicesProducer, rabbitmq *configs.RabbitMqConfig, jobCache *cache.JobCache) *JobService {
	return &JobService{
		JobRepository: jobRepo,
		Producer:      producer,
		RabbitMq:      rabbitmq,
		JobCache:      jobCache,
	}
}

func (jobService *JobService) GetJobsService(writer http.ResponseWriter, reader *http.Request) error {

	queryParams := reader.URL.Query()
	searchQuery := ""
	var provinceQuery, cityQuery int64
	if val, ok := queryParams["search"]; ok && len(val) > 0 {
		searchQuery = val[0]
	}
	if val, ok := queryParams["province_id"]; ok && len(val) > 0 {
		provinceQuery, _ = strconv.ParseInt(val[0], 10, 64)
	}
	if val, ok := queryParams["city_id"]; ok && len(val) > 0 {
		cityQuery, _ = strconv.ParseInt(val[0], 10, 64)
	}

	getJobsByRedis, err := jobService.JobCache.GetJobsFromRedis(searchQuery, provinceQuery, cityQuery)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}

	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, getJobsByRedis)
}

func (jobService *JobService) CreateJobService(request *request.JobRequest) error {

	job := &entity.JobEntity{
		ID:          uuid.NewString(),
		Title:       request.Title,
		Description: request.Description,
		Category:    request.Category,
		Price:       request.Price,
		RegencyID:   request.RegencyID,
		ProvinceID:  request.ProvinceID,
		PostedBy:    request.PostedBy,
	}

	createJob, err := jobService.JobRepository.CreateJobRepository(job)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	err = jobService.JobCache.SaveJobToRedis(createJob)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "create job is failed", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, createJob)
}

func (jobService *JobService) GetJobByIDService(reader *http.Request) error {

	vars := mux.Vars(reader)
	id := vars["id"]
	job, err := jobService.JobRepository.GetJobByIDRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "job not found", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, job)
}

func (jobService *JobService) UpdateJobService(reader *http.Request, request *request.JobRequest) error {

	vars := mux.Vars(reader)
	id := vars["id"]
	job := &entity.JobEntity{
		Title:       request.Title,
		Description: request.Description,
		Category:    request.Category,
		Price:       request.Price,
		RegencyID:   request.RegencyID,
		ProvinceID:  request.ProvinceID,
	}
	updateJob, err := jobService.JobRepository.UpdateJobRepository(id, job)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "update job is failed", http.StatusBadRequest)
	}
	err = jobService.JobCache.UpdateJobInRedis(updateJob)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, updateJob)
}

func (jobService *JobService) DeleteJobService(reader *http.Request) error {

	vars := mux.Vars(reader)
	id := vars["id"]
	err := jobService.JobRepository.DeleteJobRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}

	err = jobService.JobCache.DeleteJobFromRedis(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "Success delete job")
}

// Review Service

func (jobService *JobService) GetReviewService(writer http.ResponseWriter, reader *http.Request) error {

	vars := mux.Vars(reader)
	id := vars["id"]

	getReviews, err := jobService.JobRepository.GetReviewRepository(id)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}

	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, getReviews)
}

func (jobService *JobService) CreateReviewService(request *request.ReviewRequest) error {

	review := &entity.ReviewEntity{
		ID:         uuid.New(),
		JobID:      request.JobID,
		ReviewerID: request.ReviewerID,
		Rating:     request.Rating,
		Comment:    request.Comment,
	}
	createReview, err := jobService.JobRepository.CreateReviewRepository(review)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, createReview)
}

func (jobService *JobService) UpdateReviewService(reader *http.Request, request *request.ReviewRequest) error {

	vars := mux.Vars(reader)
	id, _ := vars["reviewId"]
	jobID, _ := vars["id"]
	parse, _ := uuid.Parse(jobID)
	review := &entity.ReviewEntity{
		JobID:      parse,
		ReviewerID: request.ReviewerID,
		Rating:     request.Rating,
		Comment:    request.Comment,
	}
	updateJob, err := jobService.JobRepository.UpdateReviewRepository(id, review)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, "update job is failed", http.StatusBadRequest)
	}
	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, updateJob)
}

func (jobService *JobService) DeleteReviewService(reader *http.Request) error {

	vars := mux.Vars(reader)
	jobID, _ := vars["id"]
	reviewID, _ := vars["reviewId"]
	err := jobService.JobRepository.DeleteReviewRepository(jobID, reviewID)
	if err != nil {
		return jobService.Producer.CreateMessageError(jobService.RabbitMq.Channel, err.Error(), http.StatusBadRequest)
	}

	return jobService.Producer.CreateMessageJob(jobService.RabbitMq.Channel, "Success delete job")
}
