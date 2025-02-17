package entity

import "time"

type TransactionEntity struct {
	ID           int64     `db:"id" json:"id"`
	JobID        int64     `db:"job_id" json:"job_id"`
	FreelancerID int64     `db:"freelancer_id" json:"freelancer_id"`
	CompanyID    int64     `db:"company_id" json:"company_id"`
	Amount       float64   `db:"amount" json:"amount"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
