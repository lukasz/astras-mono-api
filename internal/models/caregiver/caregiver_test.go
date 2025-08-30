package caregiver

import (
	"testing"
)

func TestCaregiverValidate(t *testing.T) {
	tests := []struct {
		name        string
		caregiver   Caregiver
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid caregiver",
			caregiver: Caregiver{
				Name:         "John Smith",
				Email:        "john.smith@example.com",
				Relationship: "parent",
			},
			expectError: false,
		},
		{
			name: "empty name",
			caregiver: Caregiver{
				Name:         "",
				Email:        "john.smith@example.com",
				Relationship: "parent",
			},
			expectError: true,
			errorMsg:    "name is required and cannot be empty",
		},
		{
			name: "invalid email",
			caregiver: Caregiver{
				Name:         "John Smith",
				Email:        "invalid-email",
				Relationship: "parent",
			},
			expectError: true,
			errorMsg:    "email format is invalid",
		},
		{
			name: "invalid relationship",
			caregiver: Caregiver{
				Name:         "John Smith",
				Email:        "john.smith@example.com",
				Relationship: "friend",
			},
			expectError: true,
			errorMsg:    "relationship must be one of: parent, guardian, grandparent, relative, caregiver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.caregiver.Validate()
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

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid email",
			email:       "test@example.com",
			expectError: false,
		},
		{
			name:        "valid email with subdomain",
			email:       "user.name@subdomain.example.com",
			expectError: false,
		},
		{
			name:        "valid email with plus",
			email:       "user+tag@example.com",
			expectError: false,
		},
		{
			name:        "empty email",
			email:       "",
			expectError: true,
			errorMsg:    "email is required and cannot be empty",
		},
		{
			name:        "invalid email format",
			email:       "invalid-email",
			expectError: true,
			errorMsg:    "email format is invalid",
		},
		{
			name:        "missing @ symbol",
			email:       "testexample.com",
			expectError: true,
			errorMsg:    "email format is invalid",
		},
		{
			name:        "missing domain",
			email:       "test@",
			expectError: true,
			errorMsg:    "email format is invalid",
		},
		{
			name:        "missing local part",
			email:       "@example.com",
			expectError: true,
			errorMsg:    "email format is invalid",
		},
		{
			name:        "invalid domain",
			email:       "test@domain",
			expectError: true,
			errorMsg:    "email format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
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

func TestValidateRelationship(t *testing.T) {
	tests := []struct {
		name         string
		relationship string
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "valid parent",
			relationship: "parent",
			expectError:  false,
		},
		{
			name:         "valid guardian uppercase",
			relationship: "GUARDIAN",
			expectError:  false,
		},
		{
			name:         "valid caregiver with spaces",
			relationship: "  caregiver  ",
			expectError:  false,
		},
		{
			name:         "empty relationship",
			relationship: "",
			expectError:  true,
			errorMsg:     "relationship is required",
		},
		{
			name:         "invalid relationship",
			relationship: "friend",
			expectError:  true,
			errorMsg:     "relationship must be one of: parent, guardian, grandparent, relative, caregiver",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRelationship(tt.relationship)
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

func TestGetValidRelationships(t *testing.T) {
	relationships := GetValidRelationships()
	expected := []string{"parent", "guardian", "grandparent", "relative", "caregiver"}
	
	if len(relationships) != len(expected) {
		t.Errorf("expected %d relationships, got %d", len(expected), len(relationships))
		return
	}
	
	for i, rel := range expected {
		if relationships[i] != rel {
			t.Errorf("expected relationship[%d] to be %q, got %q", i, rel, relationships[i])
		}
	}
}