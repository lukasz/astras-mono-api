package star

import (
	"testing"
)

func TestStarValidate(t *testing.T) {
	tests := []struct {
		name        string
		star        Star
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid star",
			star: Star{
				KidID:       1,
				Activity:    "homework",
				Stars:       5,
				Description: "Completed math homework",
			},
			expectError: false,
		},
		{
			name: "invalid kid ID",
			star: Star{
				KidID:       0,
				Activity:    "homework",
				Stars:       5,
				Description: "Completed math homework",
			},
			expectError: true,
			errorMsg:    "kid_id is required and must be positive",
		},
		{
			name: "empty activity",
			star: Star{
				KidID:       1,
				Activity:    "",
				Stars:       5,
				Description: "Completed math homework",
			},
			expectError: true,
			errorMsg:    "activity is required and cannot be empty",
		},
		{
			name: "invalid activity",
			star: Star{
				KidID:       1,
				Activity:    "gaming",
				Stars:       5,
				Description: "Playing games",
			},
			expectError: true,
			errorMsg:    "activity must be one of: homework, chores, behavior, helping, reading, exercise, creativity, learning",
		},
		{
			name: "too few stars",
			star: Star{
				KidID:       1,
				Activity:    "homework",
				Stars:       0,
				Description: "Attempted homework",
			},
			expectError: true,
			errorMsg:    "stars must be at least 1",
		},
		{
			name: "too many stars",
			star: Star{
				KidID:       1,
				Activity:    "homework",
				Stars:       15,
				Description: "Amazing homework",
			},
			expectError: true,
			errorMsg:    "stars cannot exceed 10",
		},
		{
			name: "description too long",
			star: Star{
				KidID:    1,
				Activity: "homework",
				Stars:    5,
				Description: "This is a very long description that exceeds the maximum allowed length of 500 characters. " +
					"It contains way too much text and should trigger a validation error because the business rules " +
					"specify that descriptions cannot be longer than 500 characters total. This description is " +
					"intentionally verbose to test the validation logic and ensure that our Star model properly " +
					"enforces the business constraint regarding maximum description length in the system. Adding more text " +
					"to ensure we exceed 500 characters in total length for this test case validation.",
			},
			expectError: true,
			errorMsg:    "description cannot exceed 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.star.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error message %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"apple", "banana", "cherry"}
	
	tests := []struct {
		item   string
		exists bool
	}{
		{"apple", true},
		{"banana", true},
		{"cherry", true},
		{"grape", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.item, func(t *testing.T) {
			result := containsString(slice, tt.item)
			if result != tt.exists {
				t.Errorf("containsString(%v, %q) = %v, want %v", slice, tt.item, result, tt.exists)
			}
		})
	}
}