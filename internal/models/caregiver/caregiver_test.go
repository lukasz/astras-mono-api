package caregiver

import (
	"testing"

	"github.com/lukasz/astras-mono-api/internal/models/caregiver/testdata"
)

func TestCaregiverValidate(t *testing.T) {
	fixture, err := testdata.LoadCaregiverValidationFixture("caregiver_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.CaregiverValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			caregiver := Caregiver{
				Name:         tt.Caregiver.Name,
				Email:        tt.Caregiver.Email,
				Relationship: RelationshipType(tt.Caregiver.Relationship),
			}

			err := caregiver.Validate()
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

func TestValidateEmail(t *testing.T) {
	fixture, err := testdata.LoadEmailValidationFixture("email_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.EmailValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			err := ValidateEmail(tt.Email)
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

func TestValidateRelationship(t *testing.T) {
	fixture, err := testdata.LoadRelationshipValidationFixture("relationship_validation_tests.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	for _, tt := range fixture.RelationshipValidationTests {
		t.Run(tt.Name, func(t *testing.T) {
			err := ValidateRelationship(tt.Relationship)
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