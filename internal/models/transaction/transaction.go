package transaction

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	// MaxDescriptionLength defines the maximum allowed description length
	MaxDescriptionLength = 255
	
	// MinStarsAmount defines the minimum stars amount for a transaction
	MinStarsAmount = 1
	
	// MaxStarsAmount defines the maximum stars amount for a single transaction
	MaxStarsAmount = 100
)

// TransactionType represents the type of star transaction
type TransactionType string

const (
	// TransactionTypeEarn represents earning stars (reward)
	TransactionTypeEarn TransactionType = "earn"
	
	// TransactionTypeSpend represents spending stars (redemption)
	TransactionTypeSpend TransactionType = "spend"
)

// Transaction represents a star transaction in the system
type Transaction struct {
	ID          int             `json:"id"`
	KidID       int             `json:"kid_id" validate:"required,min=1"`
	Type        TransactionType `json:"type" validate:"required,oneof=earn spend"`
	Amount      int             `json:"amount" validate:"required,min=1,max=100"`
	Description string          `json:"description" validate:"required,max=255"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates the transaction fields
func (t *Transaction) Validate() error {
	if err := validate.Struct(t); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "required":
				if fieldErr.Field() == "Amount" && fieldErr.Value() == 0 {
					return fmt.Errorf("amount must be at least 1")
				}
				return fmt.Errorf("%s is required", getFieldName(fieldErr.Field()))
			case "min":
				if fieldErr.Field() == "KidID" {
					return fmt.Errorf("kid_id must be greater than 0")
				}
				if fieldErr.Field() == "Amount" {
					return fmt.Errorf("amount must be at least 1")
				}
				return fmt.Errorf("%s must be at least %s", getFieldName(fieldErr.Field()), fieldErr.Param())
			case "max":
				if fieldErr.Field() == "Amount" {
					return fmt.Errorf("amount cannot exceed %d stars", MaxStarsAmount)
				}
				if fieldErr.Field() == "Description" {
					return fmt.Errorf("description cannot exceed %d characters", MaxDescriptionLength)
				}
				return fmt.Errorf("%s cannot exceed %s", getFieldName(fieldErr.Field()), fieldErr.Param())
			case "oneof":
				return fmt.Errorf("type must be either 'earn' or 'spend'")
			}
		}
	}

	// Additional business logic validation
	if err := ValidateTransactionType(string(t.Type)); err != nil {
		return err
	}

	return nil
}

// ValidateTransactionType validates if the transaction type is valid
func ValidateTransactionType(transactionType string) error {
	normalizedType := strings.TrimSpace(strings.ToLower(transactionType))
	switch TransactionType(normalizedType) {
	case TransactionTypeEarn, TransactionTypeSpend:
		return nil
	default:
		return fmt.Errorf("type must be either 'earn' or 'spend'")
	}
}

// ValidateAmount validates if the stars amount is within allowed range
func ValidateAmount(amount int) error {
	if amount < MinStarsAmount {
		return fmt.Errorf("amount must be at least %d", MinStarsAmount)
	}
	if amount > MaxStarsAmount {
		return fmt.Errorf("amount cannot exceed %d stars", MaxStarsAmount)
	}
	return nil
}

// GetValidTransactionTypes returns the list of valid transaction types
func GetValidTransactionTypes() []string {
	return []string{string(TransactionTypeEarn), string(TransactionTypeSpend)}
}

// IsEarnTransaction checks if the transaction is an earn type
func (t *Transaction) IsEarnTransaction() bool {
	return t.Type == TransactionTypeEarn
}

// IsSpendTransaction checks if the transaction is a spend type
func (t *Transaction) IsSpendTransaction() bool {
	return t.Type == TransactionTypeSpend
}

// getFieldName converts struct field names to user-friendly names
func getFieldName(field string) string {
	switch field {
	case "KidID":
		return "kid_id"
	case "Type":
		return "type"
	case "Amount":
		return "amount"
	case "Description":
		return "description"
	default:
		return field
	}
}