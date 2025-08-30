// Package main implements the Caregiver Service AWS Lambda function.
// This service manages caregivers/guardians in the Astras system, providing
// full CRUD operations and validation endpoints through a RESTful API interface.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lukasz/astras-mono-api/internal/handler"
	"github.com/lukasz/astras-mono-api/internal/models/caregiver"
)

// CaregiverRequest represents the payload for creating or updating a caregiver.
// Used for parsing JSON requests in POST and PUT operations.
type CaregiverRequest struct {
	Name         string `json:"name,omitempty"`         // Caregiver's name
	Email        string `json:"email,omitempty"`        // Contact email address
	Relationship string `json:"relationship,omitempty"` // Relationship to child
}

// ValidationRequest represents the payload for validation endpoints.
// Used for validating individual fields from frontend.
type ValidationRequest struct {
	Email        string `json:"email,omitempty"`        // Email to validate
	Relationship string `json:"relationship,omitempty"` // Relationship to validate
}

// ValidationResponse represents the response from validation endpoints.
type ValidationResponse struct {
	Valid   bool     `json:"valid"`             // Whether the value is valid
	Message string   `json:"message,omitempty"` // Error message if invalid
	Options []string `json:"options,omitempty"` // Valid options (for relationships)
}

// ToCaregiver converts a CaregiverRequest to a Caregiver model with generated fields.
// Sets timestamps and can accept an optional ID for updates.
func (cr *CaregiverRequest) ToCaregiver(id ...int) (*caregiver.Caregiver, error) {
	caregiverModel := &caregiver.Caregiver{
		Name:         strings.TrimSpace(cr.Name),
		Email:        strings.TrimSpace(cr.Email),
		Relationship: caregiver.RelationshipType(strings.TrimSpace(cr.Relationship)),
		CreatedAt:    time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		caregiverModel.ID = id[0]
		caregiverModel.UpdatedAt = time.Now()
	}

	if err := caregiverModel.Validate(); err != nil {
		return nil, err
	}

	return caregiverModel, nil
}

// CaregiverHandler implements the handler.Handler interface for caregiver-specific operations.
// This struct contains all the business logic for managing caregivers in the system.
type CaregiverHandler struct{}

// GetAll retrieves and returns a list of all caregivers in the system.
// Returns mock data for demonstration purposes - in production this would
// query a database or external service.
func (h *CaregiverHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	// Mock data - in production this would come from a database
	mockCaregivers := []caregiver.Caregiver{
		{
			ID:           1,
			Name:         "Sarah Johnson",
			Email:        "sarah.johnson@example.com",
			Relationship: caregiver.RelationshipParent,
			CreatedAt:    time.Now().AddDate(0, -3, 0), // 3 months ago
		},
		{
			ID:           2,
			Name:         "Mike Smith",
			Email:        "mike.smith@example.com",
			Relationship: caregiver.RelationshipGuardian,
			CreatedAt:    time.Now().AddDate(0, -1, 0), // 1 month ago
		},
	}

	return handler.Response{
		Message: "Caregivers retrieved successfully",
		Service: "caregiver-service",
		Data:    mockCaregivers,
	}, nil
}

// GetByID retrieves a specific caregiver by their unique identifier.
// Extracts the caregiver ID from the URL path parameters and returns mock caregiver data.
// In production, this would query the database for the specific caregiver record.
func (h *CaregiverHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid caregiver ID: %s", idStr)
	}

	// Mock data - in production this would come from a database lookup
	mockCaregiver := caregiver.Caregiver{
		ID:           id,
		Name:         "Sarah Johnson",
		Email:        "sarah.johnson@example.com",
		Relationship: caregiver.RelationshipParent,
		CreatedAt:    time.Now().AddDate(0, -3, 0),
	}

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %d retrieved successfully", id),
		Service: "caregiver-service",
		Data:    mockCaregiver,
	}, nil
}

// Create processes a request to add a new caregiver to the system.
// Parses the request body JSON and validates the caregiver data before creation.
// Returns the newly created caregiver data with a generated ID.
func (h *CaregiverHandler) Create(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var caregiverRequest CaregiverRequest
	// Parse and validate the incoming JSON request body
	if err := json.Unmarshal([]byte(request.Body), &caregiverRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model and validate
	caregiverModel, err := caregiverRequest.ToCaregiver()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// In production, save to database and get real ID
	caregiverModel.ID = 3 // Mock generated ID

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %s created successfully", caregiverModel.Name),
		Service: "caregiver-service",
		Data:    caregiverModel,
	}, nil
}

// Update modifies an existing caregiver's information in the system.
// Takes the caregiver ID from URL parameters and new data from request body.
// Returns the updated caregiver data after successful modification.
func (h *CaregiverHandler) Update(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid caregiver ID: %s", idStr)
	}

	var caregiverRequest CaregiverRequest
	// Parse and validate the incoming JSON update data
	if err := json.Unmarshal([]byte(request.Body), &caregiverRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model with existing ID and validate
	caregiverModel, err := caregiverRequest.ToCaregiver(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %d updated successfully", id),
		Service: "caregiver-service",
		Data:    caregiverModel,
	}, nil
}

// Delete removes a caregiver from the system by their unique identifier.
// Extracts the caregiver ID from URL parameters and performs the deletion operation.
// Returns a confirmation message upon successful removal.
func (h *CaregiverHandler) Delete(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid caregiver ID: %s", idStr)
	}

	// In production, perform database deletion here
	// For now, just return success

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %d deleted successfully", id),
		Service: "caregiver-service",
	}, nil
}

// ValidateEmail handles email validation requests from frontend.
// POST /validate/email with {"email": "test@example.com"}
func (h *CaregiverHandler) ValidateEmail(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var validationReq ValidationRequest
	if err := json.Unmarshal([]byte(request.Body), &validationReq); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	err := caregiver.ValidateEmail(validationReq.Email)
	response := ValidationResponse{
		Valid: err == nil,
	}
	
	if err != nil {
		response.Message = err.Error()
	}

	return handler.Response{
		Message: "Email validation completed",
		Service: "caregiver-service",
		Data:    response,
	}, nil
}

// ValidateRelationship handles relationship validation requests from frontend.
// POST /validate/relationship with {"relationship": "parent"}
func (h *CaregiverHandler) ValidateRelationship(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var validationReq ValidationRequest
	if err := json.Unmarshal([]byte(request.Body), &validationReq); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	err := caregiver.ValidateRelationship(validationReq.Relationship)
	response := ValidationResponse{
		Valid:   err == nil,
		Options: caregiver.GetValidRelationships(),
	}
	
	if err != nil {
		response.Message = err.Error()
	}

	return handler.Response{
		Message: "Relationship validation completed",
		Service: "caregiver-service",
		Data:    response,
	}, nil
}

// handleRequest is the main entry point for all HTTP requests to the Caregiver Service.
// It creates a CaregiverHandler instance and handles both CRUD operations and validation endpoints.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	caregiverHandler := &CaregiverHandler{}

	// Handle validation endpoints
	if strings.HasPrefix(request.Path, "/validate/") {
		var response handler.Response
		var err error
		
		switch {
		case strings.Contains(request.Path, "/validate/email"):
			response, err = caregiverHandler.ValidateEmail(ctx, request)
		case strings.Contains(request.Path, "/validate/relationship"):
			response, err = caregiverHandler.ValidateRelationship(ctx, request)
		default:
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Access-Control-Allow-Origin": "*",
				},
				Body: `{"message": "validation endpoint not found", "service": "caregiver-service"}`,
			}, nil
		}
		
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 400,
				Headers: map[string]string{
					"Content-Type": "application/json",
					"Access-Control-Allow-Origin": "*",
				},
				Body: fmt.Sprintf(`{"message": "%s", "service": "caregiver-service"}`, err.Error()),
			}, nil
		}
		
		responseJSON, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Access-Control-Allow-Origin": "*",
			},
			Body: string(responseJSON),
		}, nil
	}

	// Handle standard CRUD operations
	return handler.HandleRequest(ctx, request, caregiverHandler)
}

// main initializes and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	lambda.Start(handleRequest)
}