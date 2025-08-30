// Package main implements the Star Service AWS Lambda function.
// This service manages star rewards and achievements in the Astras system,
// tracking kids' accomplishments and providing CRUD operations for star rewards.
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

// StarRequest represents the payload for creating or updating a star reward.
// Used for parsing JSON requests in POST and PUT operations.
type StarRequest struct {
	KidID       int    `json:"kid_id,omitempty"`      // ID of the kid receiving the star
	Activity    string `json:"activity,omitempty"`    // Type of activity
	Stars       int    `json:"stars,omitempty"`       // Number of stars awarded
	Description string `json:"description,omitempty"` // Achievement description
}

// ToStar converts a StarRequest to a Star model with generated fields.
// Sets timestamps and can accept an optional ID for updates.
func (sr *StarRequest) ToStar(id ...int) (*models.Star, error) {
	star := &models.Star{
		KidID:       sr.KidID,
		Activity:    strings.TrimSpace(strings.ToLower(sr.Activity)),
		Stars:       sr.Stars,
		Description: strings.TrimSpace(sr.Description),
		CreatedAt:   time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		star.ID = id[0]
		star.UpdatedAt = time.Now()
	}

	if err := star.Validate(); err != nil {
		return nil, err
	}

	return star, nil
}

// StarHandler implements the handler.Handler interface for star reward operations.
// This struct contains all the business logic for managing star rewards in the system.
type StarHandler struct{}

// GetAll retrieves and returns a list of all star rewards in the system.
// Returns mock data for demonstration purposes - in production this would
// query a database or external service.
func (h *StarHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	// Mock data - in production this would come from a database
	mockStars := []models.Star{
		{
			ID:          1,
			KidID:       1,
			Activity:    "homework",
			Stars:       5,
			Description: "Completed math homework perfectly",
			CreatedAt:   time.Now().AddDate(0, 0, -7), // 7 days ago
		},
		{
			ID:          2,
			KidID:       2,
			Activity:    "chores",
			Stars:       3,
			Description: "Cleaned room thoroughly",
			CreatedAt:   time.Now().AddDate(0, 0, -3), // 3 days ago
		},
	}

	return handler.Response{
		Message: "Stars retrieved successfully",
		Service: "star-service",
		Data:    mockStars,
	}, nil
}

// GetByID retrieves a specific star reward by its unique identifier.
// Extracts the star ID from the URL path parameters and returns mock star data.
// In production, this would query the database for the specific star record.
func (h *StarHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid star ID: %s", idStr)
	}

	// Mock data - in production this would come from a database lookup
	mockStar := models.Star{
		ID:          id,
		KidID:       1,
		Activity:    "homework",
		Stars:       5,
		Description: "Completed math homework perfectly",
		CreatedAt:   time.Now().AddDate(0, 0, -7),
	}

	return handler.Response{
		Message: fmt.Sprintf("Star %d retrieved successfully", id),
		Service: "star-service",
		Data:    mockStar,
	}, nil
}

// Create processes a request to add a new star reward to the system.
// Parses the request body JSON and validates the star data before creation.
// Returns the newly created star reward data with a generated ID.
func (h *StarHandler) Create(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var starRequest StarRequest
	// Parse and validate the incoming JSON request body
	if err := json.Unmarshal([]byte(request.Body), &starRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model and validate
	star, err := starRequest.ToStar()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// In production, save to database and get real ID
	star.ID = 3 // Mock generated ID

	return handler.Response{
		Message: fmt.Sprintf("Star reward created successfully for activity: %s", star.Activity),
		Service: "star-service",
		Data:    star,
	}, nil
}

// Update modifies an existing star reward's information in the system.
// Takes the star ID from URL parameters and new data from request body.
// Returns the updated star reward data after successful modification.
func (h *StarHandler) Update(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid star ID: %s", idStr)
	}

	var starRequest StarRequest
	// Parse and validate the incoming JSON update data
	if err := json.Unmarshal([]byte(request.Body), &starRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model with existing ID and validate
	star, err := starRequest.ToStar(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Star %d updated successfully", id),
		Service: "star-service",
		Data:    star,
	}, nil
}

// Delete removes a star reward from the system by its unique identifier.
// Extracts the star ID from URL parameters and performs the deletion operation.
// Returns a confirmation message upon successful removal.
func (h *StarHandler) Delete(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid star ID: %s", idStr)
	}

	// In production, perform database deletion here
	// For now, just return success

	return handler.Response{
		Message: fmt.Sprintf("Star %d deleted successfully", id),
		Service: "star-service",
	}, nil
}

// handleRequest is the main entry point for all HTTP requests to the Star Service.
// It creates a StarHandler instance and delegates request processing to the shared
// handler infrastructure, which routes to appropriate CRUD methods based on HTTP method.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	starHandler := &StarHandler{}
	return handler.HandleRequest(ctx, request, starHandler)
}

// main initializes and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	lambda.Start(handleRequest)
}