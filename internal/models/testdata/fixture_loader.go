package testdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// FixtureTimeFormat defines the time format used in test fixtures (ISO 8601 with timezone)
	FixtureTimeFormat = time.RFC3339
)

// AgeTestCase represents a test case for Age() method
type AgeTestCase struct {
	Name        string    `json:"name"`
	Kid         KidData   `json:"kid"`
	AtTime      string    `json:"atTime"`
	ExpectedAge int       `json:"expectedAge"`
}

// ValidationTestCase represents a test case for Validate() method
type ValidationTestCase struct {
	Name                 string  `json:"name"`
	Kid                  KidData `json:"kid"`
	ExpectError          bool    `json:"expectError"`
	ExpectedErrorMessage string  `json:"expectedErrorMessage"`
}

// BirthdayTestCase represents a test case for birthday-related methods
type BirthdayTestCase struct {
	Name           string  `json:"name"`
	Kid            KidData `json:"kid"`
	CheckDate      string  `json:"checkDate"`
	ExpectedResult bool    `json:"expectedResult,omitempty"`
	ExpectedDays   int     `json:"expectedDays,omitempty"`
}

// FormatTestCase represents a test case for FormatBirthdate() method
type FormatTestCase struct {
	Name           string  `json:"name"`
	Kid            KidData `json:"kid"`
	ExpectedFormat string  `json:"expectedFormat"`
}

// KidData represents Kid data in test fixtures
type KidData struct {
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"`
}

// AgeTestFixture represents the structure of age test fixture
type AgeTestFixture struct {
	Description   string        `json:"description"`
	ReferenceDate string        `json:"referenceDate"`
	TestCases     []AgeTestCase `json:"testCases"`
}

// ValidationTestFixture represents the structure of validation test fixture
type ValidationTestFixture struct {
	Description      string               `json:"description"`
	ValidBirthdate   string               `json:"validBirthdate"`
	TestCases        []ValidationTestCase `json:"testCases"`
}

// BirthdayTestFixture represents the structure of birthday test fixture
type BirthdayTestFixture struct {
	Description   string `json:"description"`
	ReferenceDate string `json:"referenceDate"`
	TestCases     struct {
		IsBirthdayToday     []BirthdayTestCase `json:"isBirthdayToday"`
		DaysUntilBirthday   []BirthdayTestCase `json:"daysUntilBirthday"`
	} `json:"testCases"`
}

// FormatTestFixture represents the structure of format test fixture
type FormatTestFixture struct {
	Description string           `json:"description"`
	TestCases   []FormatTestCase `json:"testCases"`
}

// ParseTime parses time string to time.Time using the standard fixture format
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(FixtureTimeFormat, timeStr)
}

// LoadAgeTestFixture loads age test fixture from JSON file
func LoadAgeTestFixture(filename string) (*AgeTestFixture, error) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}

	var fixture AgeTestFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture %s: %w", filename, err)
	}

	return &fixture, nil
}

// LoadValidationTestFixture loads validation test fixture from JSON file
func LoadValidationTestFixture(filename string) (*ValidationTestFixture, error) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}

	var fixture ValidationTestFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture %s: %w", filename, err)
	}

	return &fixture, nil
}

// LoadBirthdayTestFixture loads birthday test fixture from JSON file
func LoadBirthdayTestFixture(filename string) (*BirthdayTestFixture, error) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}

	var fixture BirthdayTestFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture %s: %w", filename, err)
	}

	return &fixture, nil
}

// LoadFormatTestFixture loads format test fixture from JSON file
func LoadFormatTestFixture(filename string) (*FormatTestFixture, error) {
	data, err := os.ReadFile(filepath.Join("testdata", "fixtures", filename))
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file %s: %w", filename, err)
	}

	var fixture FormatTestFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fixture %s: %w", filename, err)
	}

	return &fixture, nil
}