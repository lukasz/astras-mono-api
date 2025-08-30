package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lukasz/astras-mono-api/internal/database/interfaces"
	"github.com/lukasz/astras-mono-api/internal/models/transaction"
)

// TransactionRepository implements the interfaces.TransactionRepository interface for PostgreSQL
type TransactionRepository struct {
	db *sqlx.DB
}

// Create adds a new transaction to the database and returns the transaction with generated ID
func (r *TransactionRepository) Create(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	// Validate the transaction before saving
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	query := `
		INSERT INTO transactions (kid_id, type, amount, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, t.KidID, string(t.Type), t.Amount, t.Description).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Return the created transaction with all data
	createdTransaction := &transaction.Transaction{
		ID:          id,
		KidID:       t.KidID,
		Type:        t.Type,
		Amount:      t.Amount,
		Description: t.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	return createdTransaction, nil
}

// GetByID retrieves a transaction by its unique identifier
func (r *TransactionRepository) GetByID(ctx context.Context, id int) (*transaction.Transaction, error) {
	query := `SELECT id, kid_id, type, amount, description, created_at, updated_at FROM transactions WHERE id = $1`

	var t transaction.Transaction
	var typeStr string
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.KidID, &typeStr, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	t.Type = transaction.TransactionType(typeStr)
	return &t, nil
}

// GetAll retrieves all transactions from the database
func (r *TransactionRepository) GetAll(ctx context.Context) ([]*transaction.Transaction, error) {
	query := `SELECT id, kid_id, type, amount, description, created_at, updated_at FROM transactions ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		var typeStr string
		
		err := rows.Scan(&t.ID, &t.KidID, &typeStr, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		
		t.Type = transaction.TransactionType(typeStr)
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

// Update modifies an existing transaction's information
func (r *TransactionRepository) Update(ctx context.Context, t *transaction.Transaction) (*transaction.Transaction, error) {
	// Validate the transaction before saving
	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("transaction validation failed: %w", err)
	}

	query := `
		UPDATE transactions 
		SET kid_id = $2, type = $3, amount = $4, description = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, kid_id, type, amount, description, created_at, updated_at`

	var updatedTransaction transaction.Transaction
	var typeStr string
	
	err := r.db.QueryRowContext(ctx, query, t.ID, t.KidID, string(t.Type), t.Amount, t.Description).Scan(
		&updatedTransaction.ID, &updatedTransaction.KidID, &typeStr, &updatedTransaction.Amount, 
		&updatedTransaction.Description, &updatedTransaction.CreatedAt, &updatedTransaction.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction with id %d not found", t.ID)
		}
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	updatedTransaction.Type = transaction.TransactionType(typeStr)
	return &updatedTransaction, nil
}

// Delete removes a transaction from the database
func (r *TransactionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction with id %d not found", id)
	}

	return nil
}

// GetByKidID retrieves all transactions for a specific kid
func (r *TransactionRepository) GetByKidID(ctx context.Context, kidID int) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, kid_id, type, amount, description, created_at, updated_at 
		FROM transactions 
		WHERE kid_id = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, kidID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by kid ID: %w", err)
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		var typeStr string
		
		err := rows.Scan(&t.ID, &t.KidID, &typeStr, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		
		t.Type = transaction.TransactionType(typeStr)
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

// GetByType retrieves all transactions of a specific type (earn/spend)
func (r *TransactionRepository) GetByType(ctx context.Context, transactionType transaction.TransactionType) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, kid_id, type, amount, description, created_at, updated_at 
		FROM transactions 
		WHERE type = $1 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(transactionType))
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by type: %w", err)
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		var typeStr string
		
		err := rows.Scan(&t.ID, &t.KidID, &typeStr, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		
		t.Type = transaction.TransactionType(typeStr)
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

// GetByKidIDAndType retrieves transactions for a specific kid and type
func (r *TransactionRepository) GetByKidIDAndType(ctx context.Context, kidID int, transactionType transaction.TransactionType) ([]*transaction.Transaction, error) {
	query := `
		SELECT id, kid_id, type, amount, description, created_at, updated_at 
		FROM transactions 
		WHERE kid_id = $1 AND type = $2 
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, kidID, string(transactionType))
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions by kid ID and type: %w", err)
	}
	defer rows.Close()

	var transactions []*transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		var typeStr string
		
		err := rows.Scan(&t.ID, &t.KidID, &typeStr, &t.Amount, &t.Description, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		
		t.Type = transaction.TransactionType(typeStr)
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transaction rows: %w", err)
	}

	return transactions, nil
}

// GetKidBalance calculates the current star balance for a kid
func (r *TransactionRepository) GetKidBalance(ctx context.Context, kidID int) (int, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'earn' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'spend' THEN amount ELSE 0 END), 0) as balance
		FROM transactions 
		WHERE kid_id = $1`

	var balance int
	err := r.db.QueryRowContext(ctx, query, kidID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to get kid balance: %w", err)
	}

	return balance, nil
}

// GetKidTransactionStats returns transaction statistics for a kid
func (r *TransactionRepository) GetKidTransactionStats(ctx context.Context, kidID int) (*interfaces.TransactionStats, error) {
	query := `
		SELECT 
			kid_id,
			COALESCE(SUM(CASE WHEN type = 'earn' THEN amount ELSE 0 END), 0) as total_earned,
			COALESCE(SUM(CASE WHEN type = 'spend' THEN amount ELSE 0 END), 0) as total_spent,
			COALESCE(SUM(CASE WHEN type = 'earn' THEN amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN type = 'spend' THEN amount ELSE 0 END), 0) as balance,
			COUNT(CASE WHEN type = 'earn' THEN 1 END) as earn_count,
			COUNT(CASE WHEN type = 'spend' THEN 1 END) as spend_count
		FROM transactions 
		WHERE kid_id = $1
		GROUP BY kid_id`

	var stats interfaces.TransactionStats
	err := r.db.QueryRowContext(ctx, query, kidID).Scan(
		&stats.KidID, &stats.TotalEarned, &stats.TotalSpent, 
		&stats.Balance, &stats.EarnCount, &stats.SpendCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no transactions found, return zero stats
			return &interfaces.TransactionStats{
				KidID:       kidID,
				TotalEarned: 0,
				TotalSpent:  0,
				Balance:     0,
				EarnCount:   0,
				SpendCount:  0,
			}, nil
		}
		return nil, fmt.Errorf("failed to get kid transaction stats: %w", err)
	}

	return &stats, nil
}