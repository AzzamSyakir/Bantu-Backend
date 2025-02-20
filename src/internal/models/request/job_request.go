package request

type JobRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Price       int    `json:"price"`
	RegencyID   string `json:"regency_id"`
	ProvinceID  string `json:"province_id"`
	PostedBy    string `json:"posted_by"`
}
