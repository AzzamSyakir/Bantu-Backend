package request

import (
	"time"
)

type Authorization struct {
	Id  string  `json:"id"`
	Rl  string  `json:"rl"`
	Exp float64 `json:"exp"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
type TopupRequest struct {
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}
type WithdrawRequest struct {
	Email         string `json:"email"`
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
}
type PayFreelancerRequest struct {
	PaymentMethod string `json:"payment_method"`
}
type AdminRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type JobRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	RegencyID   int64   `json:"regency_id"`
	ProvinceID  int64   `json:"province_id"`
	PostedBy    string  `json:"posted_by"`
}

type ProposalRequest struct {
	JobID         string   `json:"job_id"`
	FreelancerID  string   `json:"freelancer_id"`
	ProposalText  *string  `json:"proposal_text,omitempty"`
	ProposedPrice *float64 `json:"proposed_price,omitempty"`
	Status        string   `json:"status"`
}

type ReviewRequest struct {
	ID         string    `json:"id"`
	JobID      string    `json:"job_id"`
	ReviewerID string    `json:"reviewer_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
