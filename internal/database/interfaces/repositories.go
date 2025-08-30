// Package interfaces defines repository interfaces for data persistence layer.
// This package provides abstractions that can be implemented by different storage backends
// (PostgreSQL, in-memory, mock, etc.) following the Repository pattern.
package interfaces

import (
	"context"

	"github.com/lukasz/astras-mono-api/internal/models/caregiver"
	"github.com/lukasz/astras-mono-api/internal/models/kid"
	"github.com/lukasz/astras-mono-api/internal/models/transaction"
)

// KidRepository defines the interface for Kid data persistence operations.
// Implementations should handle database interactions, error handling, and data validation.
type KidRepository interface {
	// Create adds a new kid to the repository and returns the kid with generated ID
	Create(ctx context.Context, kid *kid.Kid) (*kid.Kid, error)
	
	// GetByID retrieves a kid by their unique identifier
	GetByID(ctx context.Context, id int) (*kid.Kid, error)
	
	// GetAll retrieves all kids from the repository
	GetAll(ctx context.Context) ([]*kid.Kid, error)
	
	// Update modifies an existing kid's information
	Update(ctx context.Context, kid *kid.Kid) (*kid.Kid, error)
	
	// Delete removes a kid from the repository
	Delete(ctx context.Context, id int) error
	
	// GetByAgeRange retrieves kids within a specific age range
	GetByAgeRange(ctx context.Context, minAge, maxAge int) ([]*kid.Kid, error)
}

// CaregiverRepository defines the interface for Caregiver data persistence operations.
// Implementations should handle database interactions, validation, and relationship management.
type CaregiverRepository interface {
	// Create adds a new caregiver to the repository and returns the caregiver with generated ID
	Create(ctx context.Context, caregiver *caregiver.Caregiver) (*caregiver.Caregiver, error)
	
	// GetByID retrieves a caregiver by their unique identifier
	GetByID(ctx context.Context, id int) (*caregiver.Caregiver, error)
	
	// GetAll retrieves all caregivers from the repository
	GetAll(ctx context.Context) ([]*caregiver.Caregiver, error)
	
	// Update modifies an existing caregiver's information
	Update(ctx context.Context, caregiver *caregiver.Caregiver) (*caregiver.Caregiver, error)
	
	// Delete removes a caregiver from the repository
	Delete(ctx context.Context, id int) error
	
	// GetByEmail retrieves a caregiver by their email address
	GetByEmail(ctx context.Context, email string) (*caregiver.Caregiver, error)
	
	// GetByRelationship retrieves all caregivers with a specific relationship type
	GetByRelationship(ctx context.Context, relationship caregiver.RelationshipType) ([]*caregiver.Caregiver, error)
}

// TransactionRepository defines the interface for Transaction data persistence operations.
// Implementations should handle star transaction operations, balance calculations, and kid relationships.
type TransactionRepository interface {
	// Create adds a new transaction to the repository and returns the transaction with generated ID
	Create(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error)
	
	// GetByID retrieves a transaction by its unique identifier
	GetByID(ctx context.Context, id int) (*transaction.Transaction, error)
	
	// GetAll retrieves all transactions from the repository
	GetAll(ctx context.Context) ([]*transaction.Transaction, error)
	
	// Update modifies an existing transaction's information
	Update(ctx context.Context, transaction *transaction.Transaction) (*transaction.Transaction, error)
	
	// Delete removes a transaction from the repository
	Delete(ctx context.Context, id int) error
	
	// GetByKidID retrieves all transactions for a specific kid
	GetByKidID(ctx context.Context, kidID int) ([]*transaction.Transaction, error)
	
	// GetByType retrieves all transactions of a specific type (earn/spend)
	GetByType(ctx context.Context, transactionType transaction.TransactionType) ([]*transaction.Transaction, error)
	
	// GetByKidIDAndType retrieves transactions for a specific kid and type
	GetByKidIDAndType(ctx context.Context, kidID int, transactionType transaction.TransactionType) ([]*transaction.Transaction, error)
	
	// GetKidBalance calculates the current star balance for a kid
	GetKidBalance(ctx context.Context, kidID int) (int, error)
	
	// GetKidTransactionStats returns transaction statistics for a kid (total earned, spent, balance)
	GetKidTransactionStats(ctx context.Context, kidID int) (*TransactionStats, error)
}

// TransactionStats represents aggregated transaction statistics for a kid
type TransactionStats struct {
	KidID        int `json:"kid_id"`
	TotalEarned  int `json:"total_earned"`
	TotalSpent   int `json:"total_spent"`
	Balance      int `json:"balance"`
	EarnCount    int `json:"earn_count"`
	SpendCount   int `json:"spend_count"`
}

// RepositoryManager provides access to all repository interfaces.
// This allows services to access multiple repositories through a single interface.
type RepositoryManager interface {
	// Kids returns the kid repository
	Kids() KidRepository
	
	// Caregivers returns the caregiver repository  
	Caregivers() CaregiverRepository
	
	// Transactions returns the transaction repository
	Transactions() TransactionRepository
	
	// Close closes all database connections and cleans up resources
	Close() error
	
	// Ping tests the database connection
	Ping(ctx context.Context) error
}