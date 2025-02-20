package repository

import (
	"bantu-backend/src/internal/entity"
	"database/sql"
	"time"

	"github.com/guregu/null"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	UserRepository := &UserRepository{}
	return UserRepository
}

func DeserializeUserRows(rows *sql.Rows) []*entity.UserEntity {
	var foundUsers []*entity.UserEntity
	currentTime := null.NewTime(time.Now(), true)
	for rows.Next() {
		foundUser := &entity.UserEntity{}
		scanErr := rows.Scan(
			&foundUser.ID,
			&foundUser.Name,
			&foundUser.Email,
			&foundUser.Password,
			&foundUser.Balance,
			&foundUser.Role,
			&foundUser.CreatedAt,
			&foundUser.UpdatedAt,
		)
		foundUser.UpdatedAt = currentTime.Time
		if scanErr != nil {
			panic(scanErr)
		}
		foundUsers = append(foundUsers, foundUser)
	}
	return foundUsers
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

func (userRepository *UserRepository) LoginUser(begin *sql.Tx, email string) (result *entity.UserEntity, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, email, password, balance, role, created_at, updated_at FROM "users" WHERE email=$1 LIMIT 1;`,
		email,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

	foundUsers := DeserializeUserRows(rows)
	if len(foundUsers) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundUsers[0]
	err = nil
	return result, err
}
