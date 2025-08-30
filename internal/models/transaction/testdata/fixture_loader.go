package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// TransactionTestCase represents a test case for Transaction.Validate() method
type TransactionTestCase struct {
	Name         string            `json:"name"`
	Transaction  TransactionData   `json:"transaction"`
	ExpectError  bool              `json:"expectError"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
}

// TypeValidationTestCase represents a test case for ValidateTransactionType() function
type TypeValidationTestCase struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ExpectError  bool   `json:"expectError"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// AmountValidationTestCase represents a test case for ValidateAmount() function
type AmountValidationTestCase struct {
	Name         string `json:"name"`
	Amount       int    `json:"amount"`
	ExpectError  bool   `json:"expectError"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// TransactionData represents test data for transaction model
type TransactionData struct {
	KidID       int    `json:"kid_id"`
	Type        string `json:"type"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

// TransactionValidationFixture represents the structure of transaction validation test fixture
type TransactionValidationFixture struct {
	TransactionValidationTests []TransactionTestCase `json:"transactionValidationTests"`
}

// TypeValidationFixture represents the structure of type validation test fixture
type TypeValidationFixture struct {
	TypeValidationTests []TypeValidationTestCase `json:"typeValidationTests"`
}

// AmountValidationFixture represents the structure of amount validation test fixture
type AmountValidationFixture struct {
	AmountValidationTests []AmountValidationTestCase `json:"amountValidationTests"`
}

// LoadTransactionValidationFixture loads transaction validation test cases from JSON file
func LoadTransactionValidationFixture(filename string) (*TransactionValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture TransactionValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}

// LoadTypeValidationFixture loads type validation test cases from JSON file
func LoadTypeValidationFixture(filename string) (*TypeValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture TypeValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}

// LoadAmountValidationFixture loads amount validation test cases from JSON file
func LoadAmountValidationFixture(filename string) (*AmountValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture AmountValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}