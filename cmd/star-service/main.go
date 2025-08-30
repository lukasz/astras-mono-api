// Package main implements the Star Service AWS Lambda function.
// This service manages star transactions in the Astras system,
// allowing kids to earn and spend stars through various activities.
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
	"github.com/lukasz/astras-mono-api/internal/models/transaction"
)

// TransactionRequest represents the payload for creating or updating a transaction.
// Used for parsing JSON requests in POST and PUT operations.
type TransactionRequest struct {
	KidID       int    `json:"kid_id,omitempty"`
	Type        string `json:"type,omitempty"`
	Amount      int    `json:"amount,omitempty"`
	Description string `json:"description,omitempty"`
}

// ValidationRequest represents requests to validation endpoints
type ValidationRequest struct {
	Type   string `json:"type,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

// ValidationResponse represents the response from validation endpoints
type ValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
}

// ToTransaction converts a TransactionRequest to a Transaction model with generated fields.
// Sets timestamps and can accept an optional ID for updates.
func (tr *TransactionRequest) ToTransaction(id ...int) (*transaction.Transaction, error) {
	transactionModel := &transaction.Transaction{
		KidID:       tr.KidID,
		Type:        transaction.TransactionType(strings.TrimSpace(strings.ToLower(tr.Type))),
		Amount:      tr.Amount,
		Description: strings.TrimSpace(tr.Description),
		CreatedAt:   time.Now(),
	}

	if len(id) > 0 && id[0] > 0 {
		transactionModel.ID = id[0]
		transactionModel.UpdatedAt = time.Now()
	}

	if err := transactionModel.Validate(); err != nil {
		return nil, err
	}

	return transactionModel, nil
}

// TransactionHandler implements the handler.Handler interface for star transaction operations.
// This struct contains all the business logic for managing star transactions in the system.
type TransactionHandler struct{
	repo interfaces.TransactionRepository
}

// NewTransactionHandler creates a new transaction handler with database repository
func NewTransactionHandler(repo interfaces.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		repo: repo,
	}
}

// GetAll retrieves and returns a list of all star transactions in the system.
// Uses the database repository to fetch all transactions from PostgreSQL.
func (h *TransactionHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	transactions, err := h.repo.GetAll(ctx)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to get all transactions: %w", err)
	}

	// Convert from []*transaction.Transaction to []transaction.Transaction for JSON response
	transactionList := make([]transaction.Transaction, len(transactions))
	for i, t := range transactions {
		transactionList[i] = *t
	}

	return handler.Response{
		Message: "Transactions retrieved successfully",
		Service: "star-service",
		Data:    transactionList,
	}, nil
}

// GetByID retrieves a specific star transaction by its unique identifier.
// Extracts the transaction ID from the URL path parameters and queries the database.
func (h *TransactionHandler) GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid transaction ID: %s", idStr)
	}

	transactionModel, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to get transaction: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Transaction %d retrieved successfully", id),
		Service: "star-service",
		Data:    *transactionModel,
	}, nil
}

// Create processes a request to add a new star transaction to the system.
// Parses the request body JSON and validates the transaction data before creation.
// Returns the newly created transaction data with a generated ID.
func (h *TransactionHandler) Create(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	var transactionRequest TransactionRequest
	// Parse and validate the incoming JSON request body
	if err := json.Unmarshal([]byte(request.Body), &transactionRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model and validate
	transactionModel, err := transactionRequest.ToTransaction()
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// Save to database
	createdTransaction, err := h.repo.Create(ctx, transactionModel)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Transaction created successfully: %s %d stars", createdTransaction.Type, createdTransaction.Amount),
		Service: "star-service",
		Data:    *createdTransaction,
	}, nil
}

// Update modifies an existing star transaction's information in the system.
// Takes the transaction ID from URL parameters and new data from request body.
// Returns the updated transaction data after successful modification.
func (h *TransactionHandler) Update(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid transaction ID: %s", idStr)
	}

	var transactionRequest TransactionRequest
	// Parse and validate the incoming JSON update data
	if err := json.Unmarshal([]byte(request.Body), &transactionRequest); err != nil {
		return handler.Response{}, fmt.Errorf("invalid JSON format: %v", err)
	}

	// Convert request to model with existing ID and validate
	transactionModel, err := transactionRequest.ToTransaction(id)
	if err != nil {
		return handler.Response{}, fmt.Errorf("validation failed: %v", err)
	}

	// Update in database
	updatedTransaction, err := h.repo.Update(ctx, transactionModel)
	if err != nil {
		return handler.Response{}, fmt.Errorf("failed to update transaction: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Transaction %d updated successfully", id),
		Service: "star-service",
		Data:    *updatedTransaction,
	}, nil
}

// Delete removes a star transaction from the system by its unique identifier.
// Extracts the transaction ID from URL parameters and performs the deletion operation.
// Returns a confirmation message upon successful removal.
func (h *TransactionHandler) Delete(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
	idStr := request.PathParameters["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return handler.Response{}, fmt.Errorf("invalid transaction ID: %s", idStr)
	}

	// Delete from database
	if err := h.repo.Delete(ctx, id); err != nil {
		return handler.Response{}, fmt.Errorf("failed to delete transaction: %w", err)
	}

	return handler.Response{
		Message: fmt.Sprintf("Transaction %d deleted successfully", id),
		Service: "star-service",
	}, nil
}

// HandleCustomRequest handles custom validation endpoints
func (h *TransactionHandler) HandleCustomRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path
	method := request.HTTPMethod

	headers := map[string]string{
		"Content-Type":                 "application/json",
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Headers": "Content-Type",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
	}

	if method == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
		}, nil
	}

	switch {
	case strings.HasSuffix(path, "/validate/type"):
		return h.handleTypeValidation(request, headers)
	case strings.HasSuffix(path, "/validate/amount"):
		return h.handleAmountValidation(request, headers)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Headers:    headers,
		Body:       `{"error": "endpoint not found"}`,
	}, nil
}

// handleTypeValidation validates transaction type
func (h *TransactionHandler) handleTypeValidation(request events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var req ValidationRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		response := ValidationResponse{
			Valid:   false,
			Message: "Invalid JSON format",
		}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	err := transaction.ValidateTransactionType(req.Type)
	response := ValidationResponse{
		Valid: err == nil,
	}
	if err != nil {
		response.Message = err.Error()
	}

	body, _ := json.Marshal(response)
	statusCode := 200
	if !response.Valid {
		statusCode = 400
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

// handleAmountValidation validates transaction amount
func (h *TransactionHandler) handleAmountValidation(request events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	var req ValidationRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		response := ValidationResponse{
			Valid:   false,
			Message: "Invalid JSON format",
		}
		body, _ := json.Marshal(response)
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    headers,
			Body:       string(body),
		}, nil
	}

	err := transaction.ValidateAmount(req.Amount)
	response := ValidationResponse{
		Valid: err == nil,
	}
	if err != nil {
		response.Message = err.Error()
	}

	body, _ := json.Marshal(response)
	statusCode := 200
	if !response.Valid {
		statusCode = 400
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers:    headers,
		Body:       string(body),
	}, nil
}

var transactionHandler *TransactionHandler

// initHandler initializes the transaction handler with database connection
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

	// Create transaction handler with repository
	transactionHandler = NewTransactionHandler(repoManager.Transactions())
	return nil
}

// handleRequest is the main entry point for all HTTP requests to the Star Service.
// It handles both CRUD operations and validation endpoints using the database-connected handler.
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	
	// Check for custom validation endpoints
	if strings.Contains(request.Path, "/validate/") {
		return transactionHandler.HandleCustomRequest(ctx, request)
	}
	
	return handler.HandleRequest(ctx, request, transactionHandler)
}

// main initializes the database connection and starts the AWS Lambda function handler.
// This function is called when the Lambda container starts up.
func main() {
	// Initialize handler with database connection
	if err := initHandler(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize star service: %v\n", err)
		os.Exit(1)
	}

	// Start Lambda handler
	lambda.Start(handleRequest)
}