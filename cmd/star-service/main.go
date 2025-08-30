package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type StarResponse struct {
	Message string `json:"message"`
	Service string `json:"service"`
	Data    any    `json:"data,omitempty"`
}

type StarRequest struct {
	KidID       int    `json:"kid_id,omitempty"`
	Activity    string `json:"activity,omitempty"`
	Stars       int    `json:"stars,omitempty"`
	Description string `json:"description,omitempty"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		return handleGetStars(ctx, request)
	case http.MethodPost:
		return handleCreateStar(ctx, request)
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

func handleGetStars(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := StarResponse{
		Message: "Stars retrieved successfully",
		Service: "star-service",
		Data: []map[string]any{
			{"id": 1, "kid_id": 1, "activity": "homework", "stars": 5, "description": "Completed math homework"},
			{"id": 2, "kid_id": 2, "activity": "chores", "stars": 3, "description": "Cleaned room"},
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

func handleCreateStar(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var starRequest StarRequest
	if err := json.Unmarshal([]byte(request.Body), &starRequest); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	response := StarResponse{
		Message: fmt.Sprintf("Star reward created successfully for activity: %s", starRequest.Activity),
		Service: "star-service",
		Data: map[string]any{
			"id":          3,
			"kid_id":      starRequest.KidID,
			"activity":    starRequest.Activity,
			"stars":       starRequest.Stars,
			"description": starRequest.Description,
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