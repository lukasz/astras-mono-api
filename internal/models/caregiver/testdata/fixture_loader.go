package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// CaregiverTestCase represents a test case for Caregiver.Validate() method
type CaregiverTestCase struct {
	Name         string            `json:"name"`
	Caregiver    CaregiverData     `json:"caregiver"`
	ExpectError  bool              `json:"expectError"`
	ErrorMessage string            `json:"errorMessage,omitempty"`
}

// EmailValidationTestCase represents a test case for ValidateEmail() function
type EmailValidationTestCase struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	ExpectError  bool   `json:"expectError"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// RelationshipValidationTestCase represents a test case for ValidateRelationship() function
type RelationshipValidationTestCase struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	ExpectError  bool   `json:"expectError"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// CaregiverData represents test data for caregiver model
type CaregiverData struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Relationship string `json:"relationship"`
}

// CaregiverValidationFixture represents the structure of caregiver validation test fixture
type CaregiverValidationFixture struct {
	CaregiverValidationTests []CaregiverTestCase `json:"caregiverValidationTests"`
}

// EmailValidationFixture represents the structure of email validation test fixture
type EmailValidationFixture struct {
	EmailValidationTests []EmailValidationTestCase `json:"emailValidationTests"`
}

// RelationshipValidationFixture represents the structure of relationship validation test fixture
type RelationshipValidationFixture struct {
	RelationshipValidationTests []RelationshipValidationTestCase `json:"relationshipValidationTests"`
}

// LoadCaregiverValidationFixture loads caregiver validation test cases from JSON file
func LoadCaregiverValidationFixture(filename string) (*CaregiverValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture CaregiverValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}

// LoadEmailValidationFixture loads email validation test cases from JSON file
func LoadEmailValidationFixture(filename string) (*EmailValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture EmailValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}

// LoadRelationshipValidationFixture loads relationship validation test cases from JSON file
func LoadRelationshipValidationFixture(filename string) (*RelationshipValidationFixture, error) {
	filepath := filepath.Join("testdata", "fixtures", filename)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var fixture RelationshipValidationFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, err
	}

	return &fixture, nil
}