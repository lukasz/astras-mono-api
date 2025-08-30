// Package caregiver provides the Caregiver model and related functionality for the Astras system.
package caregiver

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// RelationshipType represents the relationship between caregiver and child
type RelationshipType string

const (
	// RelationshipParent represents a parent relationship
	RelationshipParent RelationshipType = "parent"
	
	// RelationshipGuardian represents a legal guardian relationship
	RelationshipGuardian RelationshipType = "guardian"
	
	// RelationshipGrandparent represents a grandparent relationship
	RelationshipGrandparent RelationshipType = "grandparent"
	
	// RelationshipRelative represents other family relative relationship
	RelationshipRelative RelationshipType = "relative"
	
	// RelationshipCaregiver represents a non-family caregiver relationship
	RelationshipCaregiver RelationshipType = "caregiver"
)

// Caregiver represents a parent or guardian in the Astras system
// with contact information and relationship details.
type Caregiver struct {
	ID           int              `json:"id" db:"id"`                                       // Unique identifier
	Name         string           `json:"name" db:"name" validate:"required,min=2,max=100"`  // Full name
	Email        string           `json:"email" db:"email" validate:"required,email"`         // Contact email address
	Relationship RelationshipType `json:"relationship" db:"relationship" validate:"required,oneof=parent guardian grandparent relative caregiver"` // Relationship to child
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`                               // Record creation timestamp
	UpdatedAt    time.Time        `json:"updated_at,omitempty" db:"updated_at"`                   // Last update timestamp
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
	c.Relationship = RelationshipType(strings.TrimSpace(strings.ToLower(string(c.Relationship))))

	// Run struct validation
	if err := validate.Struct(c); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return formatValidationError(validationErrors)
		}
		return err
	}

	// Additional business logic validation
	if err := ValidateRelationship(string(c.Relationship)); err != nil {
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
	return []string{
		string(RelationshipParent),
		string(RelationshipGuardian),
		string(RelationshipGrandparent),
		string(RelationshipRelative),
		string(RelationshipCaregiver),
	}
}

// String returns the string representation of the RelationshipType
func (r RelationshipType) String() string {
	return string(r)
}

// IsValid checks if the relationship type is valid
func (r RelationshipType) IsValid() bool {
	switch r {
	case RelationshipParent, RelationshipGuardian, RelationshipGrandparent, RelationshipRelative, RelationshipCaregiver:
		return true
	default:
		return false
	}
}

// IsParent checks if the relationship is a parent
func (r RelationshipType) IsParent() bool {
	return r == RelationshipParent
}

// IsGuardian checks if the relationship is a legal guardian
func (r RelationshipType) IsGuardian() bool {
	return r == RelationshipGuardian
}

// IsFamily checks if the relationship represents a family member
func (r RelationshipType) IsFamily() bool {
	switch r {
	case RelationshipParent, RelationshipGrandparent, RelationshipRelative:
		return true
	default:
		return false
	}
}
