package request

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type JobRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Category    *string  `json:"category"`
	Price       *float64 `json:"price"`
	RegencyID   int64    `json:"regency_id"`
	ProvinceID  int64    `json:"province_id"`
	PostedBy    int64    `json:"posted_by"`
}
