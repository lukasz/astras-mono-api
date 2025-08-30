// Package models provides shared data structures for the Astras system.
package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Caregiver represents a parent or guardian in the Astras system
// with contact information and relationship details.
type Caregiver struct {
	ID           int       `json:"id"`                   // Unique identifier
	Name         string    `json:"name"`                 // Full name
	Email        string    `json:"email"`                // Contact email address
	Relationship string    `json:"relationship"`         // Relationship to child
	CreatedAt    time.Time `json:"created_at"`           // Record creation timestamp
	UpdatedAt    time.Time `json:"updated_at,omitempty"` // Last update timestamp
}

// Validate checks if the Caregiver data meets business requirements.
// Returns an error if any validation rules are violated.
func (c *Caregiver) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(c.Name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if len(c.Name) > 100 {
		return errors.New("name cannot exceed 100 characters")
	}

	email := strings.TrimSpace(c.Email)
	if email == "" {
		return errors.New("email is required and cannot be empty")
	}
	if !isValidEmail(email) {
		return errors.New("email format is invalid")
	}

	relationship := strings.TrimSpace(strings.ToLower(c.Relationship))
	validRelationships := []string{"parent", "guardian", "grandparent", "relative", "caregiver"}
	if !contains(validRelationships, relationship) {
		return fmt.Errorf("relationship must be one of: %s", strings.Join(validRelationships, ", "))
	}

	return nil
}

// isValidEmail performs basic email validation using simple regex-like logic.
// In production, consider using a more robust email validation library.
func isValidEmail(email string) bool {
	// Basic email validation - contains @ and has parts before/after
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	// Check for domain part structure
	domainParts := strings.Split(parts[1], ".")
	if len(domainParts) < 2 {
		return false
	}
	for _, part := range domainParts {
		if len(part) == 0 {
			return false
		}
	}
	return true
}

// contains checks if a slice contains a specific string value.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
