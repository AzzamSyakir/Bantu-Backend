package repository

import (
	"bantu-backend/src/configs"
	"bantu-backend/src/internal/entity"
	"errors"
	"fmt"
	"net/url"
)

type JobRepository struct {
	Db *configs.DatabaseConfig
}

func NewJobRepository(db *configs.DatabaseConfig) *JobRepository {
	return &JobRepository{Db: db}
}

func (jobRepository *JobRepository) GetJobsRepository(queryParams url.Values) ([]entity.JobEntity, error) {
	query := "SELECT id, title, description, category, price, regency_id, province_id, posted_by, created_at, updated_at FROM jobs WHERE 1=1"
	var args []interface{}
	argIndex := 1

	if search, exists := queryParams["search"]; exists {
		searchValue := fmt.Sprintf("%%%s%%", search[0])
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, searchValue, searchValue)
		argIndex += 2
	}
	if regencyID, exists := queryParams["regency_id"]; exists {
		query += fmt.Sprintf(" AND regency_id = $%d", argIndex)
		args = append(args, regencyID[0])
		argIndex++
	}
	if provinceID, exists := queryParams["province_id"]; exists {
		query += fmt.Sprintf(" AND province_id = $%d", argIndex)
		args = append(args, provinceID[0])
		argIndex++
	}

	rows, err := jobRepository.Db.DB.Connection.Query(query, args...)
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
			&job.Price,
			&job.RegencyID,
			&job.ProvinceID,
			&job.PostedBy,
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

func (jobRepository *JobRepository) CreateJobRepository(job *entity.JobEntity) (*entity.JobEntity, error) {
	query := `
		INSERT INTO jobs (id, title, description, category, price, regency_id, province_id, posted_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := jobRepository.Db.DB.Connection.Exec(
		query,
		job.ID,
		job.Title,
		job.Description,
		job.Category,
		job.Price,
		job.RegencyID,
		job.ProvinceID,
		job.PostedBy,
	)
	if err != nil {
		return nil, err
	}

	id := job.ID.String()
	result, err := jobRepository.GetJobByIDRepository(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (jobRepository *JobRepository) GetJobByIDRepository(id string) (*entity.JobEntity, error) {
	query := `
	SELECT j.*, r.*, p.*
	FROM jobs j
	JOIN provinces p ON j.province_id = p.id
	JOIN regencies r ON j.regency_id = r.id
	WHERE j.id = $1;
	`
	var job entity.JobEntity
	err := jobRepository.Db.DB.Connection.QueryRow(
		query,
		id,
	).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Category,
		&job.Price,
		&job.RegencyID,
		&job.ProvinceID,
		&job.PostedBy,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.Regency.ID,
		&job.Regency.ProvinceID,
		&job.Regency.RegencyName,
		&job.Regency.CreatedAt,
		&job.Province.ID,
		&job.Province.ProvinceName,
		&job.Province.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (jobRepository *JobRepository) UpdateJobRepository(id string, job *entity.JobEntity) (*entity.JobEntity, error) {
	query := `
		UPDATE jobs SET title = $1, description = $2, category = $3, regency_id = $4, province_id = $5, price = $6 WHERE id = $7
		RETURNING id, title, description, category, regency_id, province_id, price, posted_by, created_at, updated_at;
	`
	err := jobRepository.Db.DB.Connection.QueryRow(
		query,
		job.Title,
		job.Description,
		job.Category,
		job.RegencyID,
		job.ProvinceID,
		job.Price,
		id,
	).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Category,
		&job.RegencyID,
		&job.ProvinceID,
		&job.Price,
		&job.PostedBy,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (jobRepository *JobRepository) DeleteJobRepository(id string) error {
	query := `
		DELETE FROM jobs WHERE id = $1;
	`
	_, err := jobRepository.Db.DB.Connection.Exec(
		query,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (jobRepository *JobRepository) GetProposalsRepository(id string) (*[]entity.ProposalEntity, error) {
	query := `SELECT * FROM proposals WHERE job_id = $1;
	`
	rows, err := jobRepository.Db.DB.Connection.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proposals []entity.ProposalEntity
	for rows.Next() {
		var proposal entity.ProposalEntity
		err := rows.Scan(
			&proposal.ID,
			&proposal.JobID,
			&proposal.FreelancerID,
			&proposal.ProposalText,
			&proposal.ProposedPrice,
			&proposal.Status,
			&proposal.CreatedAt,
			&proposal.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		proposals = append(proposals, proposal)
	}

	return &proposals, nil
}

func (jobRepository *JobRepository) CreateProposalRepository(proposal *entity.ProposalEntity) (*entity.ProposalEntity, error) {
	query := `
		INSERT INTO proposals (id, job_id, freelancer_id, proposal_text, proposed_price, status) 
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, job_id, freelancer_id, proposal_text, proposed_price, status;
	`
	err := jobRepository.Db.DB.Connection.QueryRow(
		query,
		proposal.ID,
		proposal.JobID,
		proposal.FreelancerID,
		proposal.ProposalText,
		proposal.ProposedPrice,
		proposal.Status,
	).Scan(
		&proposal.ID,
		&proposal.JobID,
		&proposal.FreelancerID,
		&proposal.ProposalText,
		&proposal.ProposedPrice,
		&proposal.Status,
	)
	if err != nil {
		return nil, err
	}
	return proposal, nil
}

func (jobRepository *JobRepository) UpdateProposalRepository(id string, proposal *entity.ProposalEntity) (*entity.ProposalEntity, error) {
	query := `
		UPDATE proposals SET proposal_text = $1, proposed_price = $2, status = $3 WHERE id = $4
		RETURNING id, job_id, freelancer_id, proposal_text, proposed_price, status;
	`
	err := jobRepository.Db.DB.Connection.QueryRow(
		query,
		proposal.JobID,
		proposal.FreelancerID,
		proposal.ProposalText,
		proposal.ProposedPrice,
		proposal.Status,
	).Scan(
		&proposal.ID,
		&proposal.JobID,
		&proposal.FreelancerID,
		&proposal.ProposalText,
		&proposal.ProposedPrice,
		&proposal.Status,
	)
	if err != nil {
		return nil, err
	}
	return proposal, nil
}

func (jobRepository *JobRepository) AcceptProposalRepository(id string) (*entity.ProposalEntity, error) {
	query := `
		UPDATE proposals SET status = $1 WHERE id = $2 
		RETURNING id, proposal_text, proposed_price, status;
	`
	result, err := jobRepository.Db.DB.Connection.Exec(
		query,
		"accepted",
		id,
	)

	if err != nil {
		return nil, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, errors.New("proposal not found")
	}
	return nil, nil
}
