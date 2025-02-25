package entity

import "time"

type ProvinceEntity struct {
	ID           int64     `db:"id" json:"id"`
	ProvinceName string    `db:"province_name" json:"province_name"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
