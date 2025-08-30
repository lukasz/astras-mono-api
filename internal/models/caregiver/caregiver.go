// Package caregiver provides the Caregiver model and related functionality for the Astras system.
package caregiver

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Caregiver represents a parent or guardian in the Astras system
// with contact information and relationship details.
type Caregiver struct {
	ID           int       `json:"id"`                                                     // Unique identifier
	Name         string    `json:"name" validate:"required,min=2,max=100"`               // Full name
	Email        string    `json:"email" validate:"required,email"`                       // Contact email address
	Relationship string    `json:"relationship" validate:"required,oneof=parent guardian grandparent relative caregiver"` // Relationship to child
	CreatedAt    time.Time `json:"created_at"`                                             // Record creation timestamp
	UpdatedAt    time.Time `json:"updated_at,omitempty"`                                 // Last update timestamp
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate checks if the Caregiver data meets business requirements.
// Uses go-playground/validator for robust validation including proper email validation.
func (c *Caregiver) Validate() error {
	// Trim whitespace before validation
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	c.Relationship = strings.TrimSpace(strings.ToLower(c.Relationship))

	// Run struct validation
	if err := validate.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(validationErrors)
		}
		return err
	}

	return nil
}

// formatValidationError converts validator errors to user-friendly messages
func formatValidationError(errs validator.ValidationErrors) error {
	for _, err := range errs {
		switch err.Field() {
		case "Name":
			switch err.Tag() {
			case "required":
				return errors.New("name is required and cannot be empty")
			case "min":
				return errors.New("name must be at least 2 characters long")
			case "max":
				return errors.New("name cannot exceed 100 characters")
			}
		case "Email":
			switch err.Tag() {
			case "required":
				return errors.New("email is required and cannot be empty")
			case "email":
				return errors.New("email format is invalid")
			}
		case "Relationship":
			switch err.Tag() {
			case "required":
				return errors.New("relationship is required")
			case "oneof":
				return errors.New("relationship must be one of: parent, guardian, grandparent, relative, caregiver")
			}
		}
	}
	return errors.New("validation failed")
}

// ValidateEmail validates a single email address using the same validation logic
// as the Caregiver model. This can be used by frontend or other services.
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required and cannot be empty")
	}
	if err := validate.Var(email, "email"); err != nil {
		return errors.New("email format is invalid")
	}
	return nil
}

// ValidateRelationship validates a relationship value using the same logic
// as the Caregiver model. This can be used by frontend or other services.
func ValidateRelationship(relationship string) error {
	relationship = strings.TrimSpace(strings.ToLower(relationship))
	if relationship == "" {
		return errors.New("relationship is required")
	}
	if err := validate.Var(relationship, "oneof=parent guardian grandparent relative caregiver"); err != nil {
		return errors.New("relationship must be one of: parent, guardian, grandparent, relative, caregiver")
	}
	return nil
}

// GetValidRelationships returns the list of valid relationship values
func GetValidRelationships() []string {
	return []string{"parent", "guardian", "grandparent", "relative", "caregiver"}
}
