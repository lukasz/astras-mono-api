package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type KidResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Data    any    `json:"data,omitempty"`
}

type KidRequest struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		return handleGetKids(ctx, request)
	case http.MethodPost:
		return handleCreateKid(ctx, request)
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

func handleGetKids(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := KidResponse{
		Message: "Kids retrieved successfully",
		Service: "kid-service",
		Data: []map[string]any{
			{"id": 1, "name": "Alice", "age": 8},
			{"id": 2, "name": "Bob", "age": 10},
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

func handleCreateKid(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var kidRequest KidRequest
	if err := json.Unmarshal([]byte(request.Body), &kidRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	response := KidResponse{
		Message: fmt.Sprintf("Kid %s created successfully", kidRequest.Name),
		Service: "kid-service",
		Data: map[string]any{
			"id":   3,
			"name": kidRequest.Name,
			"age":  kidRequest.Age,
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