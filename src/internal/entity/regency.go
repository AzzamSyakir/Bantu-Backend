package entity

import (
	"time"
)

type RegencyEntity struct {
	ID          int64     `db:"id" json:"id"`
	ProvinceID  int64     `db:"province_id" json:"province_id"`
	RegencyName string    `db:"regency_name" json:"regency_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
