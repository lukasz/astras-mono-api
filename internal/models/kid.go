// Package models provides shared data structures for the Astras system.
package models

import (
	"errors"
	"strings"
	"time"
)

// Kid represents a child in the Astras system with personal information
// and validation rules for data integrity.
type Kid struct {
	ID        int       `json:"id"`                   // Unique identifier
	Name      string    `json:"name"`                 // Full name of the child
	Age       int       `json:"age"`                  // Age in years
	CreatedAt time.Time `json:"created_at"`           // Record creation timestamp
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Last update timestamp
}

// Validate checks if the Kid data meets business requirements.
// Returns an error if any validation rules are violated.
func (k *Kid) Validate() error {
	if strings.TrimSpace(k.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(k.Name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if len(k.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}
	if k.Age < 0 {
		return errors.New("age cannot be negative")
	}
	if k.Age > 18 {
		return errors.New("age cannot exceed 18 for kids")
	}
	return nil
}