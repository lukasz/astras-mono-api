package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lukasz/astras-mono-api/internal/models/caregiver"
)

// CaregiverRepository implements the interfaces.CaregiverRepository interface for PostgreSQL
type CaregiverRepository struct {
	db *sqlx.DB
}

// Create adds a new caregiver to the database and returns the caregiver with generated ID
func (r *CaregiverRepository) Create(ctx context.Context, c *caregiver.Caregiver) (*caregiver.Caregiver, error) {
	// Validate the caregiver before saving
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("caregiver validation failed: %w", err)
	}

	query := `
		INSERT INTO caregivers (name, email, relationship, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	var id int
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, c.Name, c.Email, string(c.Relationship)).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create caregiver: %w", err)
	}

	// Return the created caregiver with all data
	createdCaregiver := &caregiver.Caregiver{
		ID:           id,
		Name:         c.Name,
		Email:        c.Email,
		Relationship: c.Relationship,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	return createdCaregiver, nil
}

// GetByID retrieves a caregiver by their unique identifier
func (r *CaregiverRepository) GetByID(ctx context.Context, id int) (*caregiver.Caregiver, error) {
	query := `SELECT id, name, email, relationship, created_at, updated_at FROM caregivers WHERE id = $1`

	var c caregiver.Caregiver
	var relationshipStr string
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Email, &relationshipStr, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("caregiver with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get caregiver: %w", err)
	}

	c.Relationship = caregiver.RelationshipType(relationshipStr)
	return &c, nil
}

// GetAll retrieves all caregivers from the database
func (r *CaregiverRepository) GetAll(ctx context.Context) ([]*caregiver.Caregiver, error) {
	query := `SELECT id, name, email, relationship, created_at, updated_at FROM caregivers ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all caregivers: %w", err)
	}
	defer rows.Close()

	var caregivers []*caregiver.Caregiver
	for rows.Next() {
		var c caregiver.Caregiver
		var relationshipStr string
		
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &relationshipStr, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan caregiver: %w", err)
		}
		
		c.Relationship = caregiver.RelationshipType(relationshipStr)
		caregivers = append(caregivers, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating caregiver rows: %w", err)
	}

	return caregivers, nil
}

// Update modifies an existing caregiver's information
func (r *CaregiverRepository) Update(ctx context.Context, c *caregiver.Caregiver) (*caregiver.Caregiver, error) {
	// Validate the caregiver before saving
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("caregiver validation failed: %w", err)
	}

	query := `
		UPDATE caregivers 
		SET name = $2, email = $3, relationship = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, relationship, created_at, updated_at`

	var updatedCaregiver caregiver.Caregiver
	var relationshipStr string
	
	err := r.db.QueryRowContext(ctx, query, c.ID, c.Name, c.Email, string(c.Relationship)).Scan(
		&updatedCaregiver.ID, &updatedCaregiver.Name, &updatedCaregiver.Email, 
		&relationshipStr, &updatedCaregiver.CreatedAt, &updatedCaregiver.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("caregiver with id %d not found", c.ID)
		}
		return nil, fmt.Errorf("failed to update caregiver: %w", err)
	}

	updatedCaregiver.Relationship = caregiver.RelationshipType(relationshipStr)
	return &updatedCaregiver, nil
}

// Delete removes a caregiver from the database
func (r *CaregiverRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM caregivers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete caregiver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("caregiver with id %d not found", id)
	}

	return nil
}

// GetByEmail retrieves a caregiver by their email address
func (r *CaregiverRepository) GetByEmail(ctx context.Context, email string) (*caregiver.Caregiver, error) {
	query := `SELECT id, name, email, relationship, created_at, updated_at FROM caregivers WHERE email = $1`

	var c caregiver.Caregiver
	var relationshipStr string
	
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&c.ID, &c.Name, &c.Email, &relationshipStr, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("caregiver with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get caregiver by email: %w", err)
	}

	c.Relationship = caregiver.RelationshipType(relationshipStr)
	return &c, nil
}

// GetByRelationship retrieves all caregivers with a specific relationship type
func (r *CaregiverRepository) GetByRelationship(ctx context.Context, relationship caregiver.RelationshipType) ([]*caregiver.Caregiver, error) {
	query := `
		SELECT id, name, email, relationship, created_at, updated_at 
		FROM caregivers 
		WHERE relationship = $1 
		ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query, string(relationship))
	if err != nil {
		return nil, fmt.Errorf("failed to get caregivers by relationship: %w", err)
	}
	defer rows.Close()

	var caregivers []*caregiver.Caregiver
	for rows.Next() {
		var c caregiver.Caregiver
		var relationshipStr string
		
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &relationshipStr, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan caregiver: %w", err)
		}
		
		c.Relationship = caregiver.RelationshipType(relationshipStr)
		caregivers = append(caregivers, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating caregiver rows: %w", err)
	}

	return caregivers, nil
}