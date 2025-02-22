package repository

import (
	"bantu-backend/src/internal/entity"
	"database/sql"
	"time"

	"github.com/guregu/null"
)

type AdminRepository struct {
}

func NewAdminRepository() *AdminRepository {
	adminRepository := &AdminRepository{}
	return adminRepository
}

func DeserializeAdminRows(rows *sql.Rows) []*entity.AdminEntity {
	var foundAdmins []*entity.AdminEntity
	currentTime := null.NewTime(time.Now(), true)
	for rows.Next() {
		foundAdmin := &entity.AdminEntity{}
		scanErr := rows.Scan(
			&foundAdmin.ID,
			&foundAdmin.Username,
			&foundAdmin.Email,
			&foundAdmin.Password,
			&foundAdmin.CreatedAt,
			&foundAdmin.UpdatedAt,
		)
		foundAdmin.UpdatedAt = currentTime.Time
		if scanErr != nil {
			panic(scanErr)
		}
		foundAdmins = append(foundAdmins, foundAdmin)
	}
	return foundAdmins
}

func (adminRepository *AdminRepository) RegisterAdmin(begin *sql.Tx, adminEntity *entity.AdminEntity) (result *entity.AdminEntity, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "admins" (id, username, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6);`,
		adminEntity.ID,
		adminEntity.Username,
		adminEntity.Email,
		adminEntity.Password,
		adminEntity.CreatedAt,
		adminEntity.UpdatedAt,
	)
	if queryErr != nil {
		return nil, queryErr
	}

	return adminEntity, nil
}

func (adminRepository *AdminRepository) LoginAdmin(begin *sql.Tx, email string) (result *entity.AdminEntity, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, username, email, password, created_at, updated_at FROM "admins" WHERE email=$1 LIMIT 1;`,
		email,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

	foundAdmins := DeserializeAdminRows(rows)
	if len(foundAdmins) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundAdmins[0]
	err = nil
	return result, err
}
