package repository

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"log"
)

type JobRepository struct {
	Db *configs.DatabaseConfig
}

func NewJobRepository(db *configs.DatabaseConfig) *JobRepository {
	return &JobRepository{Db: db}
}

func (jobRepository *JobRepository) GetJobsRepository(db *configs.DatabaseConfig) ([]entity.JobEntity, error) {
	rows, err := db.DB.Connection.Query("SELECT id, title, description, category, location, price, posted_by, created_at, updated_at FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []entity.JobEntity

	for rows.Next() {
		var job entity.JobEntity
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Category,
			&job.Location,
			&job.Price,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}
func (jobRepository *JobRepository) CreateJobRepository(job *entity.JobEntity) (entity.JobEntity, error) {
	query := `
		INSERT INTO jobs (title, description, category, location, price, posted_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at;
	`
	err := jobRepository.Db.DB.Connection.QueryRow(
		query,
		job.Title,
		job.Description,
		job.Category,
		job.Location,
		job.Price,
		job.PostedBy,
	).Scan(
		&job.ID,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to insert job: %v", err)
		return entity.JobEntity{}, err
	}
	return *job, nil
}

func (jobRepository *JobRepository) GetJobByIDRepository(id string) (entity.JobEntity, error) {
	query := `
	SELECT id, title, description, category, location, price, posted_by, created_at, updated_at 
	FROM jobs 
	WHERE id = $1;
	`

	var job entity.JobEntity
	err := jobRepository.Db.DB.Connection.QueryRow(query, id).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Category,
		&job.Location,
		&job.Price,
		&job.PostedBy,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update job: %v", err)
		return entity.JobEntity{}, err
	}
	return job, nil
}

func (jobRepository *JobRepository) UpdateJobRepository(id string, job *entity.JobEntity) (entity.JobEntity, error) {
	query := `
		UPDATE jobs SET title = $1, description = $2, category = $3, location = $4, price = $5 WHERE id = $6;
	`
	_, err := jobRepository.Db.DB.Connection.Exec(
		query,
		job.Title,
		job.Description,
		job.Category,
		job.Location,
		job.Price,
		id,
	)
	if err != nil {
		log.Printf("Failed to update job: %v", err)
		return entity.JobEntity{}, err
	}
	return *job, nil
}

func (jobRepository *JobRepository) DeleteJobRepository(id string) error {
	query := `
		DELETE FROM jobs WHERE id = $1 
	`
	_, err := jobRepository.Db.DB.Connection.Exec(
		query,
		id,
	)
	if err != nil {
		log.Printf("Failed to update job: %v", err)
		return err
	}
	return nil
}
