package repository

import (
	"bantu-backend/src/internal/entity"
	"database/sql"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	UserRepository := &UserRepository{}
	return UserRepository
}

func (userRepository *UserRepository) RegisterUser(begin *sql.Tx, userEntitiy *entity.UserEntity) (result *entity.UserEntity, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "users" (id, name, email, password, balance, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		userEntitiy.ID,
		userEntitiy.Name,
		userEntitiy.Email,
		userEntitiy.Password,
		userEntitiy.Balance,
		userEntitiy.Role,
		userEntitiy.CreatedAt,
		userEntitiy.UpdatedAt,
	)
	if queryErr != nil {
		return nil, queryErr
	}

	return userEntitiy, nil
}
