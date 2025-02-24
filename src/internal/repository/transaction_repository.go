package repository

import (
	"bantu-backend/src/internal/entity"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}
func (transactionRepository *TransactionRepository) CreateTransaction(begin *sql.Tx, transactionEntitiy *entity.TransactionEntity) (result *entity.TransactionEntity, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO transactions (id, user_id, job_id, transaction_type, amount, payment_method, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		transactionEntitiy.ID,
		transactionEntitiy.UserId,
		transactionEntitiy.JobId,
		transactionEntitiy.TransactionType,
		transactionEntitiy.Amount,
		transactionEntitiy.PaymentMethod,
		transactionEntitiy.Status,
		transactionEntitiy.CreatedAt,
		transactionEntitiy.UpdatedAt,
	)
	if queryErr != nil {
		return nil, queryErr
	}

	return transactionEntitiy, nil
}
func (transactionRepository *TransactionRepository) UpdateTransactionStatus(tx *sql.Tx, transactionEntity *entity.TransactionEntity) error {
	fmt.Println("update transaction mulai")
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := tx.Exec(query, transactionEntity.Status, transactionEntity.UpdatedAt, transactionEntity.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("update aman")
	return nil
}
