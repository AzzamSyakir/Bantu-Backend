package entity

import (
	"time"

	"github.com/guregu/null"
)

type TransactionEntity struct {
	ID              string      `db:"id" json:"id"`
	JobId           null.String `db:"job_id" json:"job_id"`
	ProposalId      null.String `db:"proposal_id" json:"proposal_id"`
	SenderId        null.String `db:"sender_id" json:"sender_id"`
	ReceiverId      null.String `db:"receiver_id" json:"receiver_id"`
	TransactionType string      `db:"transaction_type" json:"transaction_type"`
	Amount          int         `db:"amount" json:"amount"`
	PaymentMethod   string      `db:"payment_method" json:"payment_method"`
	Status          string      `db:"status" json:"status"`
	CreatedAt       time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at" json:"updated_at"`
	InvoiceUrl      null.String `json:"invoice_url"`
	PayoutUrl       null.String `json:"payout_url"`
}
type ProposalTransactionEntity struct {
	ID            string    `db:"id" json:"id"`
	JobID         string    `db:"job_id" json:"job_id"`
	FreelancerID  string    `db:"freelancer_id" json:"freelancer_id"`
	ProposalText  *string   `db:"proposal_text" json:"proposal_text,omitempty"`
	ProposedPrice *float64  `db:"proposed_price" json:"proposed_price,omitempty"`
	Status        string    `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
