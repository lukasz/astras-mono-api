package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CaregiverResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Data    any    `json:"data,omitempty"`
}

type CaregiverRequest struct {
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Relationship string `json:"relationship,omitempty"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		return handleGetCaregivers(ctx, request)
	case http.MethodPost:
		return handleCreateCaregiver(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       `{"error": "Method not allowed"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}
}

func handleGetCaregivers(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := CaregiverResponse{
		Message: "Caregivers retrieved successfully",
		Service: "caregiver-service",
		Data: []map[string]any{
			{"id": 1, "name": "John Smith", "email": "john@example.com", "relationship": "parent"},
			{"id": 2, "name": "Jane Doe", "email": "jane@example.com", "relationship": "guardian"},
		},
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal server error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func handleCreateCaregiver(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var caregiverRequest CaregiverRequest
	if err := json.Unmarshal([]byte(request.Body), &caregiverRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	response := CaregiverResponse{
		Message: fmt.Sprintf("Caregiver %s created successfully", caregiverRequest.Name),
		Service: "caregiver-service",
		Data: map[string]any{
			"id":           3,
			"name":         caregiverRequest.Name,
			"email":        caregiverRequest.Email,
			"relationship": caregiverRequest.Relationship,
		},
	}

	body, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Internal server error"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}