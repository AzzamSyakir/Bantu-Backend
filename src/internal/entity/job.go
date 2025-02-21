package entity

import "time"

type JobEntity struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description,omitempty"`
	Category    string    `db:"category" json:"category,omitempty"`
	Price       float64   `db:"price" json:"price,omitempty"`
	RegencyID   int64     `db:"regency_id" json:"regency_id,omitempty"`
	ProvinceID  int64     `db:"province_id" json:"province_id,omitempty"`
	PostedBy    int64     `db:"posted_by" json:"posted_by"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
