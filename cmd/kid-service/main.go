// Package main implements the Kid Service AWS Lambda function.
// This service manages children/kids in the Astras system, providing
// full CRUD operations through a RESTful API interface.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/lukasz/astras-mono-api/internal/database"
	"github.com/lukasz/astras-mono-api/internal/database/interfaces"
	"github.com/lukasz/astras-mono-api/internal/database/postgres"
	"github.com/lukasz/astras-mono-api/internal/handler"
	"github.com/lukasz/astras-mono-api/internal/models/kid"
)

// KidRequest represents the payload for creating or updating a kid.
// Used for parsing JSON requests in POST and PUT operations.
type KidRequest struct {
	Name string `json:"name,omitempty"` // Child's name
	Age  int    `json:"age,omitempty"`  // Child's age
}

// ToKid converts a KidRequest to a Kid model with generated fields.
// Converts age to approximate birthdate and sets timestamps.
// Accepts an optional ID for updates.
func (kr *KidRequest) ToKid(id ...int) (*kid.Kid, error) {
	// Convert age to approximate birthdate (assuming birthday hasn't occurred this year)
	now := time.Now()
	birthdate := time.Date(now.Year()-kr.Age, now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	kidModel := &kid.Kid{
		Name:      strings.TrimSpace(kr.Name),
		Birthdate: birthdate,
		CreatedAt: time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		kidModel.ID = id[0]
		kidModel.UpdatedAt = time.Now()
	}

	if err := kidModel.Validate(); err != nil {
		return nil, err
	}

	return kidModel, nil
}

// KidHandler implements the handler.Handler interface for kid-specific operations.
// This struct contains all the business logic for managing kids in the system.
type KidHandler struct {
	repo interfaces.KidRepository
}

// NewKidHandler creates a new kid handler with database repository
func NewKidHandler(repo interfaces.KidRepository) *KidHandler {
	return &KidHandler{
		repo: repo,
	}
}

// GetAll retrieves and returns a list of all kids in the system.
// Uses the database repository to fetch all kids from PostgreSQL.
func (h *KidHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	kids, err := h.repo.GetAll(ctx)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to get all kids: %w", err)
	}

	// Convert from []*kid.Kid to []kid.Kid for JSON response
	kidList := make([]kid.Kid, len(kids))
	for i, k := range kids {
		kidList[i] = *k
	}

	return handler.Response{
		Message: "Kids retrieved successfully",
		Service: "kid-service",
		Data:    kidList,
	}, nil
}

// GetByID retrieves a specific kid by their unique identifier.
// Extracts the kid ID from the URL path parameters and queries the database.
func (h *KidHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid kid ID: %s", idStr)
	}

	kidModel, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to get kid: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %d retrieved successfully", id),
		Service: "kid-service",
		Data:    *kidModel,
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
	kidModel, err := kidRequest.ToKid()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// Save to database
	createdKid, err := h.repo.Create(ctx, kidModel)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to create kid: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %s created successfully", createdKid.Name),
		Service: "kid-service",
		Data:    *createdKid,
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
	kidModel, err := kidRequest.ToKid(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// Update in database
	updatedKid, err := h.repo.Update(ctx, kidModel)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to update kid: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %d updated successfully", id),
		Service: "kid-service",
		Data:    *updatedKid,
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

	// Delete from database
	if err := h.repo.Delete(ctx, id); err != nil {
		return handler.Response{}, fmt.Errorf("failed to delete kid: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Kid %d deleted successfully", id),
		Service: "kid-service",
	}, nil
}

var kidHandler *KidHandler

// initHandler initializes the kid handler with database connection
func initHandler() error {
	// Load database configuration from environment variables
	config := database.LoadConfigFromEnv()

	// Create PostgreSQL repository manager
	repoManager, err := postgres.NewRepositoryManager(&postgres.Config{
		Host:         config.Host,
		Port:         config.Port,
		Database:     config.Database,
		Username:     config.Username,
		Password:     config.Password,
		SSLMode:      config.SSLMode,
		MaxOpenConns: config.MaxOpenConns,
		MaxIdleConns: config.MaxIdleConns,
		MaxLifetime:  config.MaxLifetime,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Test database connection
	if err := repoManager.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Create kid handler with repository
	kidHandler = NewKidHandler(repoManager.Kids())
	return nil
}

// handleRequest is the main entry point for all HTTP requests to the Kid Service.
// It delegates request processing to the kid handler with database connectivity.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handler.HandleRequest(ctx, request, kidHandler)
}

// main initializes the database connection and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	// Initialize handler with database connection
	if err := initHandler(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize kid service: %v\n", err)
		os.Exit(1)
	}

	// Start Lambda handler
	lambda.Start(handleRequest)
}