package request

import (
	"github.com/guregu/null"
)

type RegisterRequest struct {
	Name     null.String `json:"name"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
	Role     null.String `json:"role"`
}
