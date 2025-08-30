// Package main implements the Kid Service AWS Lambda function.
// This service manages children/kids in the Astras system, providing
// full CRUD operations through a RESTful API interface.
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

// KidRequest represents the payload for creating or updating a kid.
// Used for parsing JSON requests in POST and PUT operations.
type KidRequest struct {
	Name string `json:"name,omitempty"` // Child's name
	Age  int    `json:"age,omitempty"`  // Child's age in years
}

// ToKid converts a KidRequest to a Kid model with generated fields.
// Sets timestamps and can accept an optional ID for updates.
func (kr *KidRequest) ToKid(id ...int) (*models.Kid, error) {
	kid := &models.Kid{
		Name:      strings.TrimSpace(kr.Name),
		Age:       kr.Age,
		CreatedAt: time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		kid.ID = id[0]
		kid.UpdatedAt = time.Now()
	}

	if err := kid.Validate(); err != nil {
		return nil, err
	}

	return kid, nil
}

// KidHandler implements the handler.Handler interface for kid-specific operations.
// This struct contains all the business logic for managing kids in the system.
type KidHandler struct{}

// GetAll retrieves and returns a list of all kids in the system.
// Returns mock data for demonstration purposes - in production this would
// query a database or external service.
func (h *KidHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	// Mock data - in production this would come from a database
	mockKids := []models.Kid{
		{
			ID:        1,
			Name:      "Alice Johnson",
			Age:       8,
			CreatedAt: time.Now().AddDate(0, -2, 0), // 2 months ago
		},
		{
			ID:        2,
			Name:      "Bob Smith",
			Age:       10,
			CreatedAt: time.Now().AddDate(0, -1, 0), // 1 month ago
		},
	}

	return handler.Response{
		Message: "Kids retrieved successfully",
		Service: "kid-service",
		Data:    mockKids,
	}, nil
}

// GetByID retrieves a specific kid by their unique identifier.
// Extracts the kid ID from the URL path parameters and returns mock kid data.
// In production, this would query the database for the specific kid record.
func (h *KidHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid kid ID: %s", idStr)
	}

	// Mock data - in production this would come from a database lookup
	mockKid := models.Kid{
		ID:        id,
		Name:      "Alice Johnson",
		Age:       8,
		CreatedAt: time.Now().AddDate(0, -2, 0),
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %d retrieved successfully", id),
		Service: "kid-service",
		Data:    mockKid,
	}, nil
}

// Create processes a request to add a new kid to the system.
// Parses the request body JSON and validates the kid data before creation.
// Returns the newly created kid data with a generated ID.
func (h *KidHandler) Create(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var kidRequest KidRequest
	// Parse and validate the incoming JSON request body
	if err := json.Unmarshal([]byte(request.Body), &kidRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model and validate
	kid, err := kidRequest.ToKid()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// In production, save to database and get real ID
	kid.ID = 3 // Mock generated ID

	return handler.Response{
		Message: fmt.Sprintf("Kid %s created successfully", kid.Name),
		Service: "kid-service",
		Data:    kid,
	}, nil
}

// Update modifies an existing kid's information in the system.
// Takes the kid ID from URL parameters and new data from request body.
// Returns the updated kid data after successful modification.
func (h *KidHandler) Update(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid kid ID: %s", idStr)
	}

	var kidRequest KidRequest
	// Parse and validate the incoming JSON update data
	if err := json.Unmarshal([]byte(request.Body), &kidRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model with existing ID and validate
	kid, err := kidRequest.ToKid(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %d updated successfully", id),
		Service: "kid-service",
		Data:    kid,
	}, nil
}

// Delete removes a kid from the system by their unique identifier.
// Extracts the kid ID from URL parameters and performs the deletion operation.
// Returns a confirmation message upon successful removal.
func (h *KidHandler) Delete(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid kid ID: %s", idStr)
	}

	// In production, perform database deletion here
	// For now, just return success

	return handler.Response{
		Message: fmt.Sprintf("Kid %d deleted successfully", id),
		Service: "kid-service",
	}, nil
}

// handleRequest is the main entry point for all HTTP requests to the Kid Service.
// It creates a KidHandler instance and delegates request processing to the shared
// handler infrastructure, which routes to appropriate CRUD methods based on HTTP method.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	kidHandler := &KidHandler{}
	return handler.HandleRequest(ctx, request, kidHandler)
}

// main initializes and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	lambda.Start(handleRequest)
}