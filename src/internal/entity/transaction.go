package entity

import (
	"time"

	"github.com/guregu/null"
)

type TransactionEntity struct {
	ID              string      `db:"id" json:"id"`
	UserId          string      `db:"user_id" json:"user_id"`
	JobId           null.String `db:"job_id" json:"job_id"`
	TransactionType string      `db:"transaction_type" json:"transaction_type"`
	Amount          int         `db:"amount" json:"amount"`
	PaymentMethod   string      `db:"payment_method" json:"payment_method"`
	Status          string      `db:"status" json:"status"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
	InvoiceUrl      string      `json:"invoice_url"`
}
