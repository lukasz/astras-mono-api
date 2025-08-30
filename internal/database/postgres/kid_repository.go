package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lukasz/astras-mono-api/internal/models/kid"
)

// KidRepository implements the interfaces.KidRepository interface for PostgreSQL
type KidRepository struct {
	db *sqlx.DB
}

// Create adds a new kid to the database and returns the kid with generated ID
func (r *KidRepository) Create(ctx context.Context, k *kid.Kid) (*kid.Kid, error) {
	// Validate the kid before saving
	if err := k.Validate(); err != nil {
		return nil, fmt.Errorf("kid validation failed: %w", err)
	}

	query := `
		INSERT INTO kids (name, birthdate, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, k.Name, k.Birthdate).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create kid: %w", err)
	}

	// Return the created kid with all data
	createdKid := &kid.Kid{
		ID:        id,
		Name:      k.Name,
		Birthdate: k.Birthdate,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return createdKid, nil
}

// GetByID retrieves a kid by their unique identifier
func (r *KidRepository) GetByID(ctx context.Context, id int) (*kid.Kid, error) {
	query := `SELECT id, name, birthdate, created_at, updated_at FROM kids WHERE id = $1`

	var k kid.Kid
	err := r.db.GetContext(ctx, &k, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kid with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get kid: %w", err)
	}

	return &k, nil
}

// GetAll retrieves all kids from the database
func (r *KidRepository) GetAll(ctx context.Context) ([]*kid.Kid, error) {
	query := `SELECT id, name, birthdate, created_at, updated_at FROM kids ORDER BY created_at DESC`

	var kids []kid.Kid
	err := r.db.SelectContext(ctx, &kids, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all kids: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*kid.Kid, len(kids))
	for i := range kids {
		result[i] = &kids[i]
	}

	return result, nil
}

// Update modifies an existing kid's information
func (r *KidRepository) Update(ctx context.Context, k *kid.Kid) (*kid.Kid, error) {
	// Validate the kid before saving
	if err := k.Validate(); err != nil {
		return nil, fmt.Errorf("kid validation failed: %w", err)
	}

	query := `
		UPDATE kids 
		SET name = $2, birthdate = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, birthdate, created_at, updated_at`

	var updatedKid kid.Kid
	err := r.db.QueryRowxContext(ctx, query, k.ID, k.Name, k.Birthdate).StructScan(&updatedKid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kid with id %d not found", k.ID)
		}
		return nil, fmt.Errorf("failed to update kid: %w", err)
	}

	return &updatedKid, nil
}

// Delete removes a kid from the database
func (r *KidRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM kids WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete kid: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("kid with id %d not found", id)
	}

	return nil
}

// GetByAgeRange retrieves kids within a specific age range based on birthdate calculations
func (r *KidRepository) GetByAgeRange(ctx context.Context, minAge, maxAge int) ([]*kid.Kid, error) {
	// Calculate birthdate ranges from age requirements
	now := time.Now()
	maxBirthdate := now.AddDate(-minAge, 0, 0)    // Youngest possible birthdate
	minBirthdate := now.AddDate(-maxAge-1, 0, 0)  // Oldest possible birthdate (accounting for not having birthday yet)

	query := `
		SELECT id, name, birthdate, created_at, updated_at 
		FROM kids 
		WHERE birthdate <= $1 AND birthdate > $2
		ORDER BY birthdate DESC, name ASC`

	var kids []kid.Kid
	err := r.db.SelectContext(ctx, &kids, query, maxBirthdate, minBirthdate)
	if err != nil {
		return nil, fmt.Errorf("failed to get kids by age range: %w", err)
	}

	// Convert to slice of pointers
	result := make([]*kid.Kid, len(kids))
	for i := range kids {
		result[i] = &kids[i]
	}

	return result, nil
}