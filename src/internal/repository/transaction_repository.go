package repository

import (
	"bantu-backend/src/internal/entity"
	"database/sql"
)

type TransactionRepository struct {
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}
func (transactionRepository *TransactionRepository) CreateTransaction(begin *sql.Tx, transactionEntitiy *entity.TransactionEntity) (result *entity.TransactionEntity, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO transactions (id, job_id, proposal_id, sender_id, receiver_id, transaction_type, amount, payment_method, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`,
		transactionEntitiy.ID,
		transactionEntitiy.JobId,
		transactionEntitiy.ProposalId,
		transactionEntitiy.SenderId,
		transactionEntitiy.ReceiverId,
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
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := tx.Exec(query, transactionEntity.Status, transactionEntity.UpdatedAt, transactionEntity.ID)
	if err != nil {
		panic(err)
	}
	return nil
}
