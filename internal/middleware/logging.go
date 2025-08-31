package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/lukasz/astras-mono-api/internal/logger"
)

// LoggingMiddleware provides HTTP request/response logging for Lambda functions
type LoggingMiddleware struct {
	logger      *logger.Logger
	serviceName string
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(serviceName string) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger.New(logger.Config{
			ServiceName: serviceName,
			MinLevel:    logger.INFO,
		}),
		serviceName: serviceName,
	}
}

// HandlerFunc represents a Lambda handler function
type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// WrapHandler wraps a Lambda handler with logging middleware
func (lm *LoggingMiddleware) WrapHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		startTime := time.Now()
		requestID := getRequestID(request)

		// Log incoming request
		lm.logIncomingRequest(ctx, request, requestID)

		// Execute handler
		response, err := handler(ctx, request)

		duration := time.Since(startTime)
		statusCode := response.StatusCode
		if statusCode == 0 && err != nil {
			statusCode = 500
		}

		// Log response
		lm.logResponse(ctx, request, response, err, duration, requestID, statusCode)

		return response, err
	}
}

// logIncomingRequest logs details about the incoming HTTP request
func (lm *LoggingMiddleware) logIncomingRequest(ctx context.Context, request events.APIGatewayProxyRequest, requestID string) {
	// Don't log sensitive headers
	safeHeaders := filterSensitiveHeaders(request.Headers)

	fields := []logger.Field{
		logger.RequestID(requestID),
		logger.String("http_method", request.HTTPMethod),
		logger.String("http_path", request.Path),
		logger.Any("query_params", request.QueryStringParameters),
		logger.Any("path_params", request.PathParameters),
		logger.Any("headers", safeHeaders),
		logger.String("source_ip", getSourceIP(request)),
		logger.String("user_agent", request.Headers["User-Agent"]),
	}

	// Log request body for POST/PUT requests (but limit size and filter sensitive data)
	if shouldLogRequestBody(request) {
		if body := filterSensitiveRequestBody(request.Body); body != "" {
			fields = append(fields, logger.String("request_body", body))
		}
	}

	lm.logger.Info(ctx, "Incoming HTTP request", fields...)
}

// logResponse logs details about the HTTP response
func (lm *LoggingMiddleware) logResponse(ctx context.Context, request events.APIGatewayProxyRequest, response events.APIGatewayProxyResponse, err error, duration time.Duration, requestID string, statusCode int) {
	level := logger.INFO
	message := "HTTP request completed"

	if err != nil {
		level = logger.ERROR
		message = "HTTP request failed"
	} else if statusCode >= 400 {
		level = logger.WARN
		message = "HTTP request completed with error status"
	}

	fields := []logger.Field{
		logger.RequestID(requestID),
		logger.String("http_method", request.HTTPMethod),
		logger.String("http_path", request.Path),
		logger.Int("status_code", statusCode),
		logger.Duration(duration),
	}

	// Add error information
	if err != nil {
		fields = append(fields, logger.Error(err))
	}

	// Log response body for error responses (but limit size)
	if statusCode >= 400 && response.Body != "" {
		if len(response.Body) > 1000 {
			fields = append(fields, logger.String("response_body", response.Body[:1000]+"..."))
		} else {
			fields = append(fields, logger.String("response_body", response.Body))
		}
	}

	// Add performance classification
	if duration > 5*time.Second {
		fields = append(fields, logger.String("performance", "very_slow"))
	} else if duration > 2*time.Second {
		fields = append(fields, logger.String("performance", "slow"))
	} else if duration > 1*time.Second {
		fields = append(fields, logger.String("performance", "acceptable"))
	} else {
		fields = append(fields, logger.String("performance", "fast"))
	}

	switch level {
	case logger.ERROR:
		lm.logger.Error(ctx, message, fields...)
	case logger.WARN:
		lm.logger.Warn(ctx, message, fields...)
	default:
		lm.logger.Info(ctx, message, fields...)
	}
}

// Helper functions

func getRequestID(request events.APIGatewayProxyRequest) string {
	// Try to get request ID from various headers
	if id := request.Headers["X-Request-ID"]; id != "" {
		return id
	}
	if id := request.Headers["X-Amzn-Trace-Id"]; id != "" {
		return id
	}
	if id := request.RequestContext.RequestID; id != "" {
		return id
	}
	return "unknown"
}

func getSourceIP(request events.APIGatewayProxyRequest) string {
	// Try X-Forwarded-For first (for requests through ALB/CloudFront)
	if xff := request.Headers["X-Forwarded-For"]; xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Fall back to X-Real-IP
	if realIP := request.Headers["X-Real-IP"]; realIP != "" {
		return realIP
	}

	// Finally, use the source IP from request context
	if request.RequestContext.Identity.SourceIP != "" {
		return request.RequestContext.Identity.SourceIP
	}

	return "unknown"
}

func shouldLogRequestBody(request events.APIGatewayProxyRequest) bool {
	// Only log body for POST, PUT, PATCH requests
	method := strings.ToUpper(request.HTTPMethod)
	return method == "POST" || method == "PUT" || method == "PATCH"
}

func filterSensitiveHeaders(headers map[string]string) map[string]string {
	sensitiveHeaders := map[string]bool{
		"authorization":  true,
		"cookie":         true,
		"x-api-key":      true,
		"x-auth-token":   true,
		"x-access-token": true,
	}

	filtered := make(map[string]string)
	for key, value := range headers {
		lowerKey := strings.ToLower(key)
		if sensitiveHeaders[lowerKey] {
			filtered[key] = "[REDACTED]"
		} else {
			filtered[key] = value
		}
	}
	return filtered
}

func filterSensitiveRequestBody(body string) string {
	if body == "" || len(body) > 5000 { // Don't log very large bodies
		return ""
	}

	// Try to parse as JSON and filter sensitive fields
	var jsonBody map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jsonBody); err == nil {
		// Filter sensitive fields
		sensitiveFields := []string{"password", "token", "secret", "key", "authorization"}
		for _, field := range sensitiveFields {
			if _, exists := jsonBody[field]; exists {
				jsonBody[field] = "[REDACTED]"
			}
		}

		if filtered, err := json.Marshal(jsonBody); err == nil {
			return string(filtered)
		}
	}

	// If not JSON or filtering failed, check for common sensitive patterns
	if containsSensitiveData(body) {
		return "[FILTERED_SENSITIVE_DATA]"
	}

	return body
}

func containsSensitiveData(body string) bool {
	lowerBody := strings.ToLower(body)
	sensitivePatterns := []string{
		"password",
		"passwd",
		"secret",
		"token",
		"authorization",
		"api_key",
		"apikey",
		"private_key",
		"credit_card",
		"ssn",
		"social_security",
	}

	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerBody, pattern) {
			return true
		}
	}
	return false
}

// Database logging middleware helper
func LogDatabaseOperation(ctx context.Context, serviceName, operation, table string, duration time.Duration, err error) {
	dbLogger := logger.NewDatabaseLogger(serviceName)

	_ = err // LogQuery handles error logging internally

	fields := []logger.Field{
		logger.String("operation", operation),
		logger.String("table", table),
		logger.Duration(duration),
	}

	if err != nil {
		fields = append(fields, logger.Error(err))
	}

	// Add performance classification for database operations
	if duration > 10*time.Second {
		fields = append(fields, logger.String("performance", "very_slow"))
	} else if duration > 5*time.Second {
		fields = append(fields, logger.String("performance", "slow"))
	} else if duration > 1*time.Second {
		fields = append(fields, logger.String("performance", "acceptable"))
	} else {
		fields = append(fields, logger.String("performance", "fast"))
	}

	dbLogger.LogQuery(ctx, "database_operation", []interface{}{operation, table}, duration, err)
}
