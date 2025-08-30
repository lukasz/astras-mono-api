// Package models provides shared data structures for the Astras system.
package models

import (
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const (
	// MinKidAge defines the minimum allowed age for kids (0 years)
	MinKidAge = 0
	// MaxKidAge defines the maximum allowed age for kids (18 years)
	MaxKidAge = 18

	// MinNameLength defines the minimum required length for kid names
	MinNameLength = 2
	// MaxNameLength defines the maximum allowed length for kid names
	MaxNameLength = 255

	// BirthdateFormat defines the date format for birthdate display (ISO 8601 date format)
	BirthdateFormat = time.DateOnly
)

// Kid represents a child in the Astras system with personal information
// and validation rules for data integrity.
type Kid struct {
	ID        int       `json:"id"`                   // Unique identifier
	Name      string    `json:"name"`                 // Full name of the child
	Birthdate time.Time `json:"birthdate"`            // Date of birth
	CreatedAt time.Time `json:"created_at"`           // Record creation timestamp
	UpdatedAt time.Time `json:"updated_at,omitempty"` // Last update timestamp
}

// Age calculates and returns the current age of the kid based on their birthdate.
// The calculation accounts for whether the birthday has occurred this year.
// If no time is provided, uses current time (time.Now()).
func (k *Kid) Age(at ...time.Time) int {
	var now time.Time
	if len(at) > 0 {
		now = at[0]
	} else {
		now = time.Now()
	}

	years := now.Year() - k.Birthdate.Year()

	// Check if birthday hasn't occurred yet this year
	if now.Month() < k.Birthdate.Month() ||
		(now.Month() == k.Birthdate.Month() && now.Day() < k.Birthdate.Day()) {
		years--
	}

	return years
}

// MarshalJSON customizes JSON serialization to include computed age field.
// This allows backwards compatibility with clients expecting an age field.
func (k *Kid) MarshalJSON() ([]byte, error) {
	type Alias Kid
	return json.Marshal(&struct {
		Age int `json:"age"`
		*Alias
	}{
		Age:   k.Age(),
		Alias: (*Alias)(k),
	})
}

// Validate checks if the Kid data meets business requirements.
// Returns an error if any validation rules are violated.
func (k *Kid) Validate() error {
	// Validate name
	if strings.TrimSpace(k.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if len(k.Name) < MinNameLength {
		return errors.New("name must be at least 2 characters long")
	}
	if len(k.Name) > MaxNameLength {
		return errors.New("name cannot exceed 255 characters")
	}

	// Validate birthdate
	if k.Birthdate.IsZero() {
		return errors.New("birthdate is required")
	}

	now := time.Now()
	if k.Birthdate.After(now) {
		return errors.New("birthdate cannot be in the future")
	}

	age := k.Age()
	if age < MinKidAge {
		return errors.New("invalid birthdate: results in negative age")
	}
	if age > MaxKidAge {
		return errors.New("age cannot exceed 18 for kids")
	}

	// Additional date range validation
	minDate := now.AddDate(-19, 0, 0) // 19 years ago (allows 18 year olds)

	if k.Birthdate.Before(minDate) {
		return errors.New("birthdate indicates age over 18")
	}

	return nil
}

// FormatBirthdate returns the birthdate formatted as YYYY-MM-DD string.
// Useful for displaying birthdate in a consistent format.
func (k *Kid) FormatBirthdate() string {
	return k.Birthdate.Format(BirthdateFormat)
}

// IsBirthdayToday checks if today is the kid's birthday.
// Useful for birthday notifications or special features.
// If no time is provided, uses current time (time.Now()).
func (k *Kid) IsBirthdayToday(at ...time.Time) bool {
	var now time.Time
	if len(at) > 0 {
		now = at[0]
	} else {
		now = time.Now()
	}
	return now.Month() == k.Birthdate.Month() && now.Day() == k.Birthdate.Day()
}

// DaysUntilBirthday calculates how many days until the next birthday.
// Returns 0 if today is the birthday.
// If no time is provided, uses current time (time.Now()).
func (k *Kid) DaysUntilBirthday(at ...time.Time) int {
	var now time.Time
	if len(at) > 0 {
		now = at[0]
	} else {
		now = time.Now()
	}

	// Normalize to start of day for accurate day calculation
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Get this year's birthday
	thisYearBirthday := time.Date(today.Year(), k.Birthdate.Month(), k.Birthdate.Day(), 0, 0, 0, 0, today.Location())

	// If birthday has passed or is today, calculate for next year
	if thisYearBirthday.Before(today) {
		nextYearBirthday := thisYearBirthday.AddDate(1, 0, 0)
		days := int(nextYearBirthday.Sub(today).Hours() / 24)
		return days
	}

	// Birthday is today or in the future
	days := int(thisYearBirthday.Sub(today).Hours() / 24)
	return days
}
