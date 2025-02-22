package request

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type JobRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Price       float64   `json:"price"`
	RegencyID   int64     `json:"regency_id"`
	ProvinceID  int64     `json:"province_id"`
	PostedBy    uuid.UUID `json:"posted_by"`
}

type ProposalRequest struct {
	JobID         uuid.UUID `json:"job_id"`
	FreelancerID  uuid.UUID `json:"freelancer_id"`
	ProposalText  *string   `json:"proposal_text,omitempty"`
	ProposedPrice *float64  `json:"proposed_price,omitempty"`
	Status        string    `json:"status"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
