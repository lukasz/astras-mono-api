# Logging Solution for Astras API

Complete logging infrastructure for the Astras microservices architecture using structured logging, CloudWatch integration, and automated monitoring.

## Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Lambda        │ -> │  Structured      │ -> │  CloudWatch     │
│   Functions     │    │  JSON Logs       │    │  Log Groups     │  
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │  Metric Filters │
                                               │  & Alarms       │
                                               └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │  CloudWatch     │
                                               │  Dashboard      │
                                               └─────────────────┘
```

## Components

### 1. Structured Logger (`internal/logger/`)
- **JSON-based logging** for easy parsing and analysis
- **Multiple log levels**: DEBUG, INFO, WARN, ERROR
- **Context-aware** with AWS Lambda integration
- **Field-based logging** with typed fields
- **Performance classification** for database operations
- **Sensitive data filtering** for security

### 2. HTTP Middleware (`internal/middleware/`)
- **Automatic request/response logging**
- **Performance monitoring** with timing
- **Error tracking** with stack traces
- **Request ID correlation** across services
- **Sensitive header filtering** for security

### 3. CloudWatch Integration (`serverless-logging.yml`)
- **Centralized log groups** with retention policies
- **Metric filters** for error detection
- **Automated alarms** for monitoring
- **CloudWatch Dashboard** for visualization

## Features

### ✅ Structured Logging
```json
{
  "@timestamp": "2024-01-15T10:30:45Z",
  "level": "INFO",
  "service": "kid-service", 
  "message": "HTTP request completed",
  "request_id": "abc-123-def",
  "http_method": "POST",
  "http_path": "/kids",
  "status_code": 201,
  "duration": 245,
  "environment": "dev",
  "aws_request_id": "lambda-req-123"
}
```

### ✅ Performance Monitoring
- **Request duration tracking**
- **Database operation timing**  
- **Performance classification**: fast/acceptable/slow/very_slow
- **Slow query detection** with automatic alerting

### ✅ Error Tracking
- **Automatic error logging** with stack traces
- **Error metrics** and alerting
- **Error rate monitoring**
- **Database error detection**

### ✅ Security Features
- **Sensitive data filtering** (passwords, tokens, keys)
- **Header sanitization** for logs
- **Request body filtering** for sensitive fields
- **No PII logging** in production

## Deployment

### 1. Deploy Logging Infrastructure
```bash
# Deploy CloudWatch resources
mage deploy:logging

# Or manually
serverless deploy --config serverless-logging.yml --stage dev
```

### 2. Deploy Services with Logging
```bash
# Deploy all services with logging enabled
mage deploy:all
```

### 3. View Dashboard
After deployment, access your dashboard:
```bash
# Get dashboard URL from CloudFormation outputs
aws cloudformation describe-stacks \
  --stack-name astras-logging-dev \
  --query 'Stacks[0].Outputs[?OutputKey==`DashboardURL`].OutputValue' \
  --output text
```

## Usage Examples

### Basic Logging
```go
import "internal/logger"

// Initialize logger
appLogger := logger.New(logger.Config{
    ServiceName: "kid-service",
    MinLevel:    logger.INFO,
})

// Log with context and fields
appLogger.Info(ctx, "User created successfully", 
    logger.String("user_id", "123"),
    logger.Int("age", 25),
    logger.Duration(processingTime),
)
```

### HTTP Middleware
```go
import "internal/middleware"

// Initialize middleware
loggingMiddleware := middleware.NewLoggingMiddleware("kid-service")

// Wrap handler
wrappedHandler := loggingMiddleware.WrapHandler(yourHandler)
```

### Database Logging
```go
import "internal/middleware"

// Log database operations automatically
middleware.LogDatabaseOperation(ctx, "kid-service", "INSERT", "kids", duration, err)
```

### Error Logging
```go
// Errors are automatically logged with stack traces
appLogger.Error(ctx, "Failed to create kid",
    logger.Error(err),
    logger.String("operation", "create_kid"),
    logger.Any("request_data", kidRequest),
)
```

## Log Retention Policies

| Environment | Retention | Rationale |
|------------|-----------|-----------|
| **dev** | 7 days | Short-term debugging |
| **staging** | 30 days | Integration testing |
| **prod** | 90 days | Compliance & analysis |

## Monitoring & Alerts

### Automatic Alarms
- **High Error Rate**: >5 errors in 5 minutes
- **Database Errors**: Any database connection failures
- **Slow Queries**: Queries taking >10 seconds

### Metrics Tracked
- **Error Count** by service
- **Database Error Count**
- **Slow Query Count**
- **Request Duration** percentiles

### Dashboard Widgets
1. **Recent Errors** - Last 100 error messages
2. **Error Metrics** - Error counts over time
3. **Slow Database Operations** - Performance issues
4. **Service Health** - Overall system status

## Log Analysis Queries

### Find Recent Errors
```sql
SOURCE '/astras/dev/application'
| fields @timestamp, level, service, message, error
| filter level = "ERROR"
| sort @timestamp desc
| limit 100
```

### Database Performance Issues  
```sql
SOURCE '/astras/dev/database'
| fields @timestamp, operation, duration, table
| filter duration > 1000
| sort duration desc
| limit 50
```

### Request Patterns by Service
```sql
SOURCE '/astras/dev/application'
| fields @timestamp, service, http_method, http_path, status_code
| stats count() by service, http_method, http_path
| sort count desc
```

### Error Rate Over Time
```sql
SOURCE '/astras/dev/application'
| fields @timestamp, level
| filter @timestamp > @timestamp - 1h
| stats count() as total_requests, 
        sum(level="ERROR") as errors
        by bin(5m)
| sort @timestamp desc
```

## Cost Estimation

### CloudWatch Logs Costs (per month)
- **Data Ingestion**: ~$0.50 per GB
- **Storage**: ~$0.03 per GB per month
- **Insights Queries**: ~$0.005 per GB scanned

### Estimated Monthly Costs
| Environment | Log Volume | Monthly Cost |
|------------|------------|--------------|
| **dev** | ~1 GB | ~$0.53 |
| **staging** | ~5 GB | ~$2.65 |  
| **prod** | ~20 GB | ~$10.60 |

## Best Practices

### ✅ Do
- **Use structured JSON logging** for all services
- **Include request IDs** for correlation
- **Log business events** (user actions, transactions)
- **Set appropriate log levels** (avoid DEBUG in production)
- **Use metric filters** for automated monitoring
- **Filter sensitive data** before logging

### ❌ Don't
- **Log sensitive data** (passwords, tokens, PII)
- **Log large payloads** without size limits
- **Use string concatenation** for log messages
- **Log at DEBUG level** in production
- **Ignore log retention costs**
- **Skip error context** information

## Troubleshooting

### High Log Volume
```bash
# Check log group sizes
aws logs describe-log-groups \
  --log-group-name-prefix "/astras" \
  --query 'logGroups[*].[logGroupName,storedBytes]' \
  --output table
```

### Missing Logs
```bash
# Verify Lambda function permissions
aws iam get-role-policy \
  --role-name astras-kid-service-dev-lambda-role \
  --policy-name CloudWatchLogsPolicy
```

### Alarm Not Firing
```bash
# Check metric filter
aws logs describe-metric-filters \
  --log-group-name "/astras/dev/application"
```

## Integration Examples

### With Existing Handler
```go
// Before (no logging)
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return handler.HandleRequest(ctx, request, kidHandler)
}

// After (with logging)
func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    wrappedHandler := loggingMiddleware.WrapHandler(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        return handler.HandleRequest(ctx, request, kidHandler)
    })
    
    return wrappedHandler(ctx, request)
}
```

### Custom Fields
```go
// Add custom business context
appLogger.Info(ctx, "Kid created successfully",
    logger.String("kid_id", kid.ID),
    logger.String("caregiver_id", caregiver.ID),
    logger.Int("kid_age", kid.Age),
    logger.String("business_event", "kid_registration"),
)
```

## Next Steps

1. **Deploy logging infrastructure** first
2. **Update all services** to use structured logging
3. **Configure alerts** for your specific thresholds
4. **Set up log analysis** queries for your use cases
5. **Monitor costs** and adjust retention policies
6. **Train team** on log analysis techniques

For questions or issues, check the CloudWatch Insights documentation or review the dashboard for real-time system health.