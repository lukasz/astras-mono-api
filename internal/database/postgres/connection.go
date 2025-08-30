// Package postgres provides PostgreSQL implementations of repository interfaces.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver (pgx)

	"github.com/lukasz/astras-mono-api/internal/database/interfaces"
)

// Config holds the PostgreSQL database configuration
type Config struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// DefaultConfig returns a default PostgreSQL configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         5432,
		Database:     "astras",
		Username:     "postgres",
		Password:     "",
		SSLMode:      "disable",
		MaxOpenConns: 25,
		MaxIdleConns: 5,
		MaxLifetime:  5 * time.Minute,
	}
}

// DSN returns the PostgreSQL data source name
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// RepositoryManager implements the interfaces.RepositoryManager interface
type RepositoryManager struct {
	db           *sqlx.DB
	kidRepo      *KidRepository
	caregiverRepo *CaregiverRepository
	transactionRepo *TransactionRepository
}

// NewRepositoryManager creates a new PostgreSQL repository manager
func NewRepositoryManager(config *Config) (*RepositoryManager, error) {
	db, err := sqlx.Connect("pgx", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.MaxLifetime)

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create repository instances
	rm := &RepositoryManager{
		db: db,
	}
	
	rm.kidRepo = &KidRepository{db: db}
	rm.caregiverRepo = &CaregiverRepository{db: db}
	rm.transactionRepo = &TransactionRepository{db: db}

	return rm, nil
}

// Kids returns the kid repository
func (rm *RepositoryManager) Kids() interfaces.KidRepository {
	return rm.kidRepo
}

// Caregivers returns the caregiver repository
func (rm *RepositoryManager) Caregivers() interfaces.CaregiverRepository {
	return rm.caregiverRepo
}

// Transactions returns the transaction repository
func (rm *RepositoryManager) Transactions() interfaces.TransactionRepository {
	return rm.transactionRepo
}

// Close closes the database connection
func (rm *RepositoryManager) Close() error {
	if rm.db != nil {
		return rm.db.Close()
	}
	return nil
}

// Ping tests the database connection
func (rm *RepositoryManager) Ping(ctx context.Context) error {
	return rm.db.PingContext(ctx)
}

// GetDB returns the underlying database connection (for testing/migrations)
func (rm *RepositoryManager) GetDB() *sqlx.DB {
	return rm.db
}

// withTx executes a function within a database transaction
func (rm *RepositoryManager) withTx(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := rm.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}