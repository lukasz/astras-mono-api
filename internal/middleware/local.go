package middleware

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"internal/logger"
)

// LocalLoggingMiddleware provides logging for local development
type LocalLoggingMiddleware struct {
	logger      *logger.Logger
	localLogger *logger.LocalLogger
	serviceName string
	isLocal     bool
}

// NewLocalLoggingMiddleware creates middleware for local development
func NewLocalLoggingMiddleware(serviceName string) (*LocalLoggingMiddleware, error) {
	// Detect if we're running locally
	isLocal := isLocalEnvironment()

	var appLogger *logger.Logger
	var localLogger *logger.LocalLogger
	var err error

	if isLocal {
		// Use local file logging for development
		localLogger, err = logger.LocalDevelopmentSetup(serviceName)
		if err != nil {
			return nil, err
		}
		appLogger = localLogger.Logger
	} else {
		// Use regular CloudWatch logging for Lambda
		appLogger = logger.New(logger.Config{
			ServiceName: serviceName,
			MinLevel:    logger.INFO,
		})
	}

	return &LocalLoggingMiddleware{
		logger:      appLogger,
		localLogger: localLogger,
		serviceName: serviceName,
		isLocal:     isLocal,
	}, nil
}

// WrapHandler wraps a handler with appropriate logging for the environment
func (lm *LocalLoggingMiddleware) WrapHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if lm.isLocal {
			return lm.wrapLocalHandler(handler)(ctx, request)
		}
		
		// Use regular middleware for Lambda
		middleware := &LoggingMiddleware{
			logger:      lm.logger,
			serviceName: lm.serviceName,
		}
		return middleware.WrapHandler(handler)(ctx, request)
	}
}

// wrapLocalHandler provides enhanced logging for local development
func (lm *LocalLoggingMiddleware) wrapLocalHandler(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		// Enhanced local logging with more details
		lm.logger.Debug(ctx, "=== LOCAL DEVELOPMENT REQUEST ===")
		lm.logger.Info(ctx, "Incoming request",
			logger.String("method", request.HTTPMethod),
			logger.String("path", request.Path),
			logger.Any("headers", request.Headers),
			logger.Any("query_params", request.QueryStringParameters),
			logger.Any("path_params", request.PathParameters),
		)

		if request.Body != "" {
			lm.logger.Debug(ctx, "Request body", logger.String("body", request.Body))
		}

		// Execute handler
		response, err := handler(ctx, request)

		// Log response
		if err != nil {
			lm.logger.Error(ctx, "Request failed", logger.Error(err))
		} else {
			lm.logger.Info(ctx, "Request completed",
				logger.Int("status_code", response.StatusCode),
				logger.String("response_body", response.Body),
			)
		}

		lm.logger.Debug(ctx, "=== END REQUEST ===")
		return response, err
	}
}

// Close closes any local resources
func (lm *LocalLoggingMiddleware) Close() error {
	if lm.localLogger != nil {
		return lm.localLogger.Close()
	}
	return nil
}

// isLocalEnvironment detects if we're running in local development
func isLocalEnvironment() bool {
	// Check various environment indicators
	stage := os.Getenv("STAGE")
	environment := os.Getenv("ENVIRONMENT")
	awsRegion := os.Getenv("AWS_REGION")
	awsLambdaFunction := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

	// If no AWS Lambda function name, we're likely local
	if awsLambdaFunction == "" {
		return true
	}

	// Check for local development stage names
	localStages := []string{"local", "dev-local", "development"}
	for _, localStage := range localStages {
		if strings.EqualFold(stage, localStage) || strings.EqualFold(environment, localStage) {
			return true
		}
	}

	// Check for SAM local indicators
	if strings.Contains(awsRegion, "local") || awsRegion == "" {
		return true
	}

	return false
}

// SetupLocalLogging is a convenience function for setting up local logging
func SetupLocalLogging(serviceName string) (*LocalLoggingMiddleware, error) {
	middleware, err := NewLocalLoggingMiddleware(serviceName)
	if err != nil {
		return nil, err
	}

	// Log setup information
	if middleware.isLocal {
		middleware.logger.Info(context.Background(), "Local logging initialized",
			logger.String("service", serviceName),
			logger.String("environment", "local"),
		)
	}

	return middleware, nil
}