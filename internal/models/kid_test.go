package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lukasz/astras-mono-api/internal/models/testdata"
)

const (
	// Test constants for consistent test values
	testKidAge  = 10              // Standard age used in non-deterministic tests
	testKidName = "Alice Johnson" // Standard name used in benchmarks and tests
)

// TestKidAge tests the Age() method using JSON fixtures
func TestKidAge(t *testing.T) {
	fixture, err := testdata.LoadAgeTestFixture("kid_age_tests.json")
	if err != nil {
		t.Fatalf("Failed to load age test fixture: %v", err)
	}

	for _, tc := range fixture.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse birthdate
			birthdate, err := testdata.ParseTime(tc.Kid.Birthdate)
			if err != nil {
				t.Fatalf("Failed to parse birthdate %s: %v", tc.Kid.Birthdate, err)
			}

			// Parse reference time
			atTime, err := testdata.ParseTime(tc.AtTime)
			if err != nil {
				t.Fatalf("Failed to parse atTime %s: %v", tc.AtTime, err)
			}

			kid := &Kid{
				Name:      tc.Kid.Name,
				Birthdate: birthdate,
			}

			gotAge := kid.Age(atTime)
			if gotAge != tc.ExpectedAge {
				t.Errorf("Age() = %d, want %d", gotAge, tc.ExpectedAge)
			}
		})
	}
}

// TestKidValidate tests the Validate() method using JSON fixtures
func TestKidValidate(t *testing.T) {
	fixture, err := testdata.LoadValidationTestFixture("kid_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load validation test fixture: %v", err)
	}

	for _, tc := range fixture.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse birthdate
			birthdate, err := testdata.ParseTime(tc.Kid.Birthdate)
			if err != nil {
				t.Fatalf("Failed to parse birthdate %s: %v", tc.Kid.Birthdate, err)
			}

			kid := &Kid{
				Name:      tc.Kid.Name,
				Birthdate: birthdate,
			}

			err = kid.Validate()

			if tc.ExpectError {
				if err == nil {
					t.Errorf("Validate() error = nil, want error containing %q", tc.ExpectedErrorMessage)
				} else if tc.ExpectedErrorMessage != "" && err.Error() != tc.ExpectedErrorMessage {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tc.ExpectedErrorMessage)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

// TestKidFormatBirthdate tests the FormatBirthdate() method using JSON fixtures
func TestKidFormatBirthdate(t *testing.T) {
	fixture, err := testdata.LoadFormatTestFixture("kid_format_tests.json")
	if err != nil {
		t.Fatalf("Failed to load format test fixture: %v", err)
	}

	for _, tc := range fixture.TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse birthdate
			birthdate, err := testdata.ParseTime(tc.Kid.Birthdate)
			if err != nil {
				t.Fatalf("Failed to parse birthdate %s: %v", tc.Kid.Birthdate, err)
			}

			kid := &Kid{
				Name:      tc.Kid.Name,
				Birthdate: birthdate,
			}

			got := kid.FormatBirthdate()
			if got != tc.ExpectedFormat {
				t.Errorf("FormatBirthdate() = %q, want %q", got, tc.ExpectedFormat)
			}
		})
	}
}

// TestKidIsBirthdayToday tests the IsBirthdayToday() method using JSON fixtures
func TestKidIsBirthdayToday(t *testing.T) {
	fixture, err := testdata.LoadBirthdayTestFixture("kid_birthday_tests.json")
	if err != nil {
		t.Fatalf("Failed to load birthday test fixture: %v", err)
	}

	for _, tc := range fixture.TestCases.IsBirthdayToday {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse birthdate
			birthdate, err := testdata.ParseTime(tc.Kid.Birthdate)
			if err != nil {
				t.Fatalf("Failed to parse birthdate %s: %v", tc.Kid.Birthdate, err)
			}

			// Parse check date
			checkDate, err := testdata.ParseTime(tc.CheckDate)
			if err != nil {
				t.Fatalf("Failed to parse checkDate %s: %v", tc.CheckDate, err)
			}

			kid := &Kid{
				Name:      tc.Kid.Name,
				Birthdate: birthdate,
			}

			got := kid.IsBirthdayToday(checkDate)
			if got != tc.ExpectedResult {
				t.Errorf("IsBirthdayToday() = %v, want %v", got, tc.ExpectedResult)
			}
		})
	}
}

// TestKidDaysUntilBirthday tests the DaysUntilBirthday() method using JSON fixtures
func TestKidDaysUntilBirthday(t *testing.T) {
	fixture, err := testdata.LoadBirthdayTestFixture("kid_birthday_tests.json")
	if err != nil {
		t.Fatalf("Failed to load birthday test fixture: %v", err)
	}

	for _, tc := range fixture.TestCases.DaysUntilBirthday {
		t.Run(tc.Name, func(t *testing.T) {
			// Parse birthdate
			birthdate, err := testdata.ParseTime(tc.Kid.Birthdate)
			if err != nil {
				t.Fatalf("Failed to parse birthdate %s: %v", tc.Kid.Birthdate, err)
			}

			// Parse check date
			checkDate, err := testdata.ParseTime(tc.CheckDate)
			if err != nil {
				t.Fatalf("Failed to parse checkDate %s: %v", tc.CheckDate, err)
			}

			kid := &Kid{
				Name:      tc.Kid.Name,
				Birthdate: birthdate,
			}

			got := kid.DaysUntilBirthday(checkDate)
			if got != tc.ExpectedDays {
				t.Errorf("DaysUntilBirthday() = %d, want %d", got, tc.ExpectedDays)
			}
		})
	}
}

// TestKidMarshalJSON tests the JSON marshaling with computed age field
func TestKidMarshalJSON(t *testing.T) {
	// Use a fixed date for consistent testing
	fixedDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	fixedBirthdate := time.Date(2016, 6, 15, 0, 0, 0, 0, time.UTC)

	kid := &Kid{
		ID:        1,
		Name:      testKidName,
		Birthdate: fixedBirthdate, // Will be 8 years old on fixedDate
		CreatedAt: fixedDate,
	}

	data, err := json.Marshal(kid)
	if err != nil {
		t.Fatalf("Failed to marshal Kid: %v", err)
	}

	// Parse JSON to check if age field is present
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check if age field exists
	age, exists := result["age"]
	if !exists {
		t.Error("JSON output missing 'age' field")
	}

	// Check if age value is correct (based on current date, not fixed date)
	// This is a limitation - MarshalJSON uses time.Now() internally
	if ageFloat, ok := age.(float64); ok {
		expectedAge := kid.Age() // Uses time.Now()
		if int(ageFloat) != expectedAge {
			t.Errorf("JSON age = %d, want %d", int(ageFloat), expectedAge)
		}
	} else {
		t.Errorf("Age field is not a number: %T", age)
	}

	// Check other required fields
	requiredFields := []string{"id", "name", "birthdate", "created_at"}
	for _, field := range requiredFields {
		if _, exists := result[field]; !exists {
			t.Errorf("JSON output missing required field: %s", field)
		}
	}
}

// TestKidAgeWithoutParameter tests that Age() works without parameter
func TestKidAgeWithoutParameter(t *testing.T) {
	// This test verifies the method works with time.Now() as default
	kid := &Kid{
		Name:      "Test Kid",
		Birthdate: time.Now().AddDate(-testKidAge, 0, 0), // testKidAge years ago from now
	}

	age := kid.Age() // Should use time.Now() internally
	if age != testKidAge {
		t.Errorf("Age() without parameter = %d, want %d", age, testKidAge)
	}
}

// BenchmarkKidAge benchmarks the Age() calculation
func BenchmarkKidAge(b *testing.B) {
	kid := &Kid{
		Name:      "Test Kid",
		Birthdate: time.Date(2014, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	referenceDate := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kid.Age(referenceDate)
	}
}

// BenchmarkKidValidate benchmarks the Validate() method
func BenchmarkKidValidate(b *testing.B) {
	kid := &Kid{
		Name:      testKidName,
		Birthdate: time.Date(2016, 6, 15, 0, 0, 0, 0, time.UTC),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = kid.Validate()
	}
}

// BenchmarkKidMarshalJSON benchmarks JSON marshaling
func BenchmarkKidMarshalJSON(b *testing.B) {
	kid := &Kid{
		ID:        1,
		Name:      testKidName,
		Birthdate: time.Date(2016, 6, 15, 0, 0, 0, 0, time.UTC),
		CreatedAt: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(kid)
	}
}
