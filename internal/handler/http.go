// Package handler provides shared HTTP request handling infrastructure for AWS Lambda functions.
// This package centralizes common HTTP operations, response formatting, and error handling
// to eliminate code duplication across microservices.
package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Response represents the standardized JSON response structure for all API endpoints.
// It provides consistent formatting across all microservices with optional data payload.
type Response struct {
	Message string `json:"message"`          // Human-readable message describing the operation result
	Service string `json:"service"`          // Name of the service that handled the request
	Data    any    `json:"data,omitempty"`   // Optional data payload (omitted if nil/empty)
}

// Handler defines the contract that all service handlers must implement.
// Each method corresponds to a standard CRUD operation and returns a Response and error.
// This interface enables polymorphic handling of different resource types.
type Handler interface {
	// GetAll retrieves and returns a list of all resources of this type
	GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error)
	
	// GetByID retrieves a specific resource by its unique identifier
	GetByID(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error)
	
	// Create processes a request to create a new resource instance
	Create(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error)
	
	// Update modifies an existing resource identified by ID with new data
	Update(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error)
	
	// Delete removes a resource identified by ID from the system
	Delete(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error)
}

// HandleRequest is the central HTTP request router and processor for AWS Lambda functions.
// It examines the HTTP method and path parameters to determine which handler method to call,
// then formats the response consistently with proper HTTP status codes and CORS headers.
// This function eliminates the need for each service to implement its own routing logic.
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest, h Handler) (events.APIGatewayProxyResponse, error) {
	var response Response
	var err error
	var statusCode int

	// Route the request to the appropriate handler method based on HTTP method and path
	switch request.HTTPMethod {
	case http.MethodGet:
		// Check if this is a single resource request (has ID parameter) or list request
		if id := request.PathParameters["id"]; id != "" {
			response, err = h.GetByID(ctx, request)
		} else {
			response, err = h.GetAll(ctx, request)
		}
		statusCode = http.StatusOK
	case http.MethodPost:
		// Handle resource creation requests
		response, err = h.Create(ctx, request)
		statusCode = http.StatusCreated  // 201 for successful creation
	case http.MethodPut:
		// Handle resource update requests
		response, err = h.Update(ctx, request)
		statusCode = http.StatusOK
	case http.MethodDelete:
		// Handle resource deletion requests
		response, err = h.Delete(ctx, request)
		statusCode = http.StatusOK
	default:
		// Return 405 Method Not Allowed for unsupported HTTP methods
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error": "Method not allowed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Handle any errors returned by the handler methods
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "` + err.Error() + `"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Marshal the response to JSON format
	body, err := json.Marshal(response)
	if err != nil {
		// Return 500 Internal Server Error if JSON marshaling fails
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal server error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	// Return successful response with appropriate status code and CORS headers
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",  // Enable CORS for all origins
		},
	}, nil
}