// Package models provides shared data structures for the Astras system.
package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Star represents a reward/achievement in the Astras system
// tracking kids' accomplishments and activities.
type Star struct {
	ID          int       `json:"id"`                   // Unique identifier
	KidID       int       `json:"kid_id"`               // ID of the kid receiving the star
	Activity    string    `json:"activity"`             // Type of activity completed
	Stars       int       `json:"stars"`                // Number of stars awarded
	Description string    `json:"description"`          // Detailed description of achievement
	CreatedAt   time.Time `json:"created_at"`           // Record creation timestamp
	UpdatedAt   time.Time `json:"updated_at,omitempty"` // Last update timestamp
}

// Validate checks if the Star data meets business requirements.
// Returns an error if any validation rules are violated.
func (s *Star) Validate() error {
	if s.KidID <= 0 {
		return errors.New("kid_id is required and must be positive")
	}

	if strings.TrimSpace(s.Activity) == "" {
		return errors.New("activity is required and cannot be empty")
	}

	activity := strings.TrimSpace(strings.ToLower(s.Activity))
	validActivities := []string{"homework", "chores", "behavior", "helping", "reading", "exercise", "creativity", "learning"}
	if !containsString(validActivities, activity) {
		return fmt.Errorf("activity must be one of: %s", strings.Join(validActivities, ", "))
	}

	if s.Stars < 1 {
		return errors.New("stars must be at least 1")
	}
	if s.Stars > 10 {
		return errors.New("stars cannot exceed 10")
	}

	if len(s.Description) > 500 {
		return errors.New("description cannot exceed 500 characters")
	}

	return nil
}

// containsString checks if a slice contains a specific string value.
func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}