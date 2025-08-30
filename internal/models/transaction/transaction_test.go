package transaction

import (
	"strings"
	"testing"

	"github.com/lukasz/astras-mono-api/internal/models/transaction/testdata"
)

func TestTransactionValidate(t *testing.T) {
	fixture, err := testdata.LoadTransactionValidationFixture("transaction_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.TransactionValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			transaction := Transaction{
				KidID:       tt.Transaction.KidID,
				Type:        TransactionType(strings.TrimSpace(strings.ToLower(tt.Transaction.Type))),
				Amount:      tt.Transaction.Amount,
				Description: strings.TrimSpace(tt.Transaction.Description),
			}

			err := transaction.Validate()
			if tt.ExpectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.ErrorMessage != "" && err.Error() != tt.ErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.ErrorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateTransactionType(t *testing.T) {
	fixture, err := testdata.LoadTypeValidationFixture("type_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.TypeValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			err := ValidateTransactionType(tt.Type)
			if tt.ExpectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.ErrorMessage != "" && err.Error() != tt.ErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.ErrorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	fixture, err := testdata.LoadAmountValidationFixture("amount_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.AmountValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			err := ValidateAmount(tt.Amount)
			if tt.ExpectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.ErrorMessage != "" && err.Error() != tt.ErrorMessage {
					t.Errorf("expected error message %q, got %q", tt.ErrorMessage, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGetValidTransactionTypes(t *testing.T) {
	types := GetValidTransactionTypes()
	expected := []string{"earn", "spend"}
	
	if len(types) != len(expected) {
		t.Errorf("expected %d transaction types, got %d", len(expected), len(types))
		return
	}
	
	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("expected type[%d] to be %q, got %q", i, expectedType, types[i])
		}
	}
}

func TestTransactionTypeMethods(t *testing.T) {
	earnTransaction := Transaction{
		KidID:       1,
		Type:        TransactionTypeEarn,
		Amount:      5,
		Description: "Test earn",
	}
	
	spendTransaction := Transaction{
		KidID:       1,
		Type:        TransactionTypeSpend,
		Amount:      3,
		Description: "Test spend",
	}
	
	// Test IsEarnTransaction
	if !earnTransaction.IsEarnTransaction() {
		t.Error("expected earn transaction to return true for IsEarnTransaction()")
	}
	
	if spendTransaction.IsEarnTransaction() {
		t.Error("expected spend transaction to return false for IsEarnTransaction()")
	}
	
	// Test IsSpendTransaction
	if earnTransaction.IsSpendTransaction() {
		t.Error("expected earn transaction to return false for IsSpendTransaction()")
	}
	
	if !spendTransaction.IsSpendTransaction() {
		t.Error("expected spend transaction to return true for IsSpendTransaction()")
	}
}