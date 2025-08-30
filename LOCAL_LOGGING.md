# Local Logging Setup for Astras API

Complete guide for setting up logging during local development, providing the same structured logging experience as production CloudWatch.

## Quick Start

### 1. Simple File-Based Logging
```bash
# Start basic logging infrastructure
./scripts/local-logging.sh start

# Build and run with SAM Local
mage build:allLocal
sam local start-api --env-vars local-env.json --docker-network astras-mono-api_astras-network --port 3000

# View logs in real-time
./scripts/local-logging.sh logs
```

### 2. Advanced ELK Stack Logging
```bash
# Start with Elasticsearch + Kibana
./scripts/local-logging.sh start --elastic

# Build and run services
mage build:allLocal
sam local start-api --env-vars local-env.json --docker-network astras-logging --port 3000

# Open Kibana dashboard
./scripts/local-logging.sh dashboard
```

## Architecture Options

### Option 1: File-Based Logging (Recommended for Development)
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   SAM Local     │ -> │  Local Files     │ -> │  Fluent Bit     │
│   Lambda        │    │  (JSON logs)     │    │  Aggregator     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │  Log Viewer     │
                                               │  (Web UI)       │
                                               └─────────────────┘
```

### Option 2: ELK Stack (Production-Like)
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   SAM Local     │ -> │  Logstash        │ -> │  Elasticsearch  │
│   Lambda        │    │  (Processing)    │    │  (Storage)      │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                        │
                                               ┌─────────────────┐
                                               │  Kibana         │
                                               │  (Dashboard)    │
                                               └─────────────────┘
```

## Features

### ✅ **Automatic Environment Detection**
```go
// Logger automatically detects local vs Lambda environment
middleware, err := middleware.NewLocalLoggingMiddleware("kid-service")
// Uses file logging locally, CloudWatch in Lambda
```

### ✅ **Enhanced Local Logging**
- **More verbose output** (DEBUG level enabled)
- **Request/response body logging** for debugging
- **File rotation** with timestamps
- **Multi-output** (console + files)

### ✅ **Same JSON Structure**
```json
{
  "@timestamp": "2024-01-15T10:30:45Z",
  "level": "DEBUG",
  "service": "kid-service",
  "message": "=== LOCAL DEVELOPMENT REQUEST ===", 
  "environment": "local",
  "http_method": "POST",
  "http_path": "/kids"
}
```

## Commands Reference

### Local Logging Script
```bash
# Start logging infrastructure
./scripts/local-logging.sh start [--elastic]

# Stop all services
./scripts/local-logging.sh stop

# View real-time logs
./scripts/local-logging.sh logs

# Open dashboard
./scripts/local-logging.sh dashboard

# Clean log data
./scripts/local-logging.sh clean

# Check service status
./scripts/local-logging.sh status
```

### Development Workflow
```bash
# 1. Start logging
./scripts/local-logging.sh start

# 2. Start database
docker-compose up -d

# 3. Build services
mage build:allLocal

# 4. Start SAM Local
sam local start-api --env-vars local-env.json --port 3000

# 5. View logs (in another terminal)
./scripts/local-logging.sh logs

# 6. Test endpoints
curl http://localhost:3000/kids
```

## File Locations

### Log Files
```
logs/
├── astras-kid-service-2024-01-15.log
├── astras-caregiver-service-2024-01-15.log  
├── astras-star-service-2024-01-15.log
├── astras-migration-service-2024-01-15.log
└── aggregated.log (combined)
```

### Configuration Files
```
local-env.json              # Environment variables for SAM Local
docker-compose.logging.yml  # Logging infrastructure
fluent-bit.conf             # Log aggregation config
logstash/                   # ELK Stack configuration
└── pipeline/
    └── logstash.conf
```

## Integration Examples

### Using in Service Code
```go
// Option 1: Auto-detection (recommended)
func initHandler() error {
    // Automatically detects local vs Lambda
    loggingMiddleware, err := middleware.NewLocalLoggingMiddleware("kid-service")
    if err != nil {
        return err
    }
    
    // Use same way as production
    wrappedHandler := loggingMiddleware.WrapHandler(yourHandler)
    return nil
}

// Option 2: Explicit local logger
func initLocalHandler() error {
    localLogger, err := logger.LocalDevelopmentSetup("kid-service")
    if err != nil {
        return err
    }
    defer localLogger.Close() // Important: close file handles
    
    localLogger.Info(ctx, "Service started locally")
    return nil
}
```

### Custom Local Configuration
```go
config := logger.LocalConfig{
    Config: logger.Config{
        ServiceName: "kid-service",
        Environment: "local",
        MinLevel:    logger.DEBUG,
    },
    LogDir:      "custom-logs",
    LogFileName: "my-service.log",
    MaxFileSize: 10 * 1024 * 1024, // 10MB
}

localLogger, err := logger.NewLocalLogger(config)
```

## ELK Stack Setup

### Kibana Dashboard Creation

1. **Open Kibana**: http://localhost:5601

2. **Create Index Pattern**:
   - Go to Management → Index Patterns
   - Create pattern: `astras-logs-*`
   - Select timestamp field: `@timestamp`

3. **Import Dashboard**:
   - Go to Management → Saved Objects
   - Import pre-configured dashboards

4. **Useful Visualizations**:
   - **Error Timeline**: Errors over time by service
   - **Request Volume**: Requests per minute
   - **Performance Metrics**: Response time percentiles
   - **Service Health**: Status codes distribution

### Sample Kibana Queries
```
# Recent errors
level: "ERROR" AND @timestamp: [now-1h TO now]

# Slow requests  
duration: >1000

# Service-specific logs
service: "kid-service" AND level: "INFO"

# Database operations
operation: EXISTS AND table: EXISTS
```

## Performance Impact

### File-Based Logging
- **Minimal impact** on local development
- **~1-2ms overhead** per request
- **Disk usage**: ~1MB per hour of active development

### ELK Stack
- **Moderate resource usage**: ~512MB RAM for Elasticsearch
- **Better for integration testing** with production-like setup
- **Rich analysis capabilities**

## Troubleshooting

### Logs Not Appearing
```bash
# Check service status
./scripts/local-logging.sh status

# Verify log directory permissions
ls -la logs/

# Check SAM Local output
sam local start-api --debug
```

### ELK Stack Issues
```bash
# Check Elasticsearch health
curl http://localhost:9200/_cluster/health

# View Logstash logs
docker logs astras-logstash

# Restart services
docker-compose -f docker-compose.logging.yml restart
```

### File Permission Issues
```bash
# Fix log directory permissions
sudo chown -R $USER:$USER logs/
chmod 755 logs/
```

## Production Parity

### What's the Same
- ✅ **JSON log format** identical to CloudWatch
- ✅ **Structured fields** for easy parsing
- ✅ **Request correlation** with IDs
- ✅ **Error tracking** and stack traces
- ✅ **Performance metrics** and timing

### What's Different  
- 📁 **File storage** instead of CloudWatch
- 🔍 **More verbose** logging (DEBUG enabled)
- 💻 **Local dashboards** instead of AWS Console
- 🔄 **Manual log rotation** (not automatic)

## Migration to Production

### Code Changes Required
**None!** The same code works in both environments:

```go
// This works locally AND in Lambda
loggingMiddleware := middleware.NewLocalLoggingMiddleware("service-name")
wrappedHandler := loggingMiddleware.WrapHandler(handler)
```

### Deployment Differences
```bash
# Local development
sam local start-api --env-vars local-env.json

# Production deployment  
mage deploy:all
```

The logging middleware automatically detects the environment and uses the appropriate logging backend.

## Best Practices

### ✅ Do
- **Start logging infrastructure** before running services
- **Use structured fields** for consistent parsing
- **Close log files** properly in cleanup
- **Monitor disk usage** in long-running sessions
- **Use DEBUG level** for detailed local debugging

### ❌ Don't  
- **Mix logging approaches** in the same codebase
- **Commit log files** to Git (add `logs/` to `.gitignore`)
- **Run ELK stack** on low-memory machines (<4GB RAM)
- **Log sensitive data** even in local development
- **Ignore log rotation** for long sessions

This setup provides production-quality logging for local development while maintaining the same developer experience as your cloud infrastructure!