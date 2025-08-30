// Package main implements the Caregiver Service AWS Lambda function.
// This service manages caregivers and guardians in the Astras system, providing
// full CRUD operations for parent/guardian profiles through a RESTful API.
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
	"github.com/lukasz/astras-mono-api/internal/models"
)

// CaregiverRequest represents the payload for creating or updating a caregiver.
// Used for parsing JSON requests in POST and PUT operations.
type CaregiverRequest struct {
	Name         string `json:"name,omitempty"`         // Caregiver's full name
	Email        string `json:"email,omitempty"`        // Contact email address
	Relationship string `json:"relationship,omitempty"` // Relationship to child
}

// ToCaregiver converts a CaregiverRequest to a Caregiver model with generated fields.
// Sets timestamps and can accept an optional ID for updates.
func (cr *CaregiverRequest) ToCaregiver(id ...int) (*models.Caregiver, error) {
	caregiver := &models.Caregiver{
		Name:         strings.TrimSpace(cr.Name),
		Email:        strings.TrimSpace(strings.ToLower(cr.Email)),
		Relationship: strings.TrimSpace(strings.ToLower(cr.Relationship)),
		CreatedAt:    time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		caregiver.ID = id[0]
		caregiver.UpdatedAt = time.Now()
	}

	if err := caregiver.Validate(); err != nil {
		return nil, err
	}

	return caregiver, nil
}

// CaregiverHandler implements the handler.Handler interface for caregiver-specific operations.
// This struct contains all the business logic for managing caregivers in the system.
type CaregiverHandler struct{}

// GetAll retrieves and returns a list of all caregivers in the system.
// Returns mock data for demonstration purposes - in production this would
// query a database or external service.
func (h *CaregiverHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	// Mock data - in production this would come from a database
	mockCaregivers := []models.Caregiver{
		{
			ID:           1,
			Name:         "John Smith",
			Email:        "john.smith@example.com",
			Relationship: "parent",
			CreatedAt:    time.Now().AddDate(0, -3, 0), // 3 months ago
		},
		{
			ID:           2,
			Name:         "Jane Doe",
			Email:        "jane.doe@example.com",
			Relationship: "guardian",
			CreatedAt:    time.Now().AddDate(0, -2, 0), // 2 months ago
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
	mockCaregiver := models.Caregiver{
		ID:           id,
		Name:         "John Smith",
		Email:        "john.smith@example.com",
		Relationship: "parent",
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
	caregiver, err := caregiverRequest.ToCaregiver()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// In production, save to database and get real ID
	caregiver.ID = 3 // Mock generated ID

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %s created successfully", caregiver.Name),
		Service: "caregiver-service",
		Data:    caregiver,
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
	caregiver, err := caregiverRequest.ToCaregiver(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Caregiver %d updated successfully", id),
		Service: "caregiver-service",
		Data:    caregiver,
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

// handleRequest is the main entry point for all HTTP requests to the Caregiver Service.
// It creates a CaregiverHandler instance and delegates request processing to the shared
// handler infrastructure, which routes to appropriate CRUD methods based on HTTP method.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	caregiverHandler := &CaregiverHandler{}
	return handler.HandleRequest(ctx, request, caregiverHandler)
}

// main initializes and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	lambda.Start(handleRequest)
}