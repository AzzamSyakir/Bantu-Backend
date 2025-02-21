package request

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type JobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	RegencyID   string `json:"regency_id"`
	ProvinceID  string `json:"province_id"`
	PostedBy    string `json:"posted_by"`
}
