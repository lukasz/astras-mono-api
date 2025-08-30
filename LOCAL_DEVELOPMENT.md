# Local Development Guide

Guide for running Lambda functions locally using AWS SAM.

## 🚀 Quick Start

### Requirements
- Go 1.21+
- AWS SAM CLI
- Mage build tool
- Docker (dla SAM CLI)

### Tool Installation

```bash
# Install AWS SAM CLI (macOS)
brew tap aws/tap
brew install aws-sam-cli

# Install Mage
go install github.com/magefile/mage@latest

# Verify installation
sam --version
mage -l
```

## 🏗️ Building and Running

### 1. Building Service
```bash
# Add Mage to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin

# Build kid-service
mage build:kid

# Or build all services
mage build:all
```

### 2. Starting Local API
```bash
# Start SAM local API on port 3000
sam local start-api --port 3000

# API will be available at: http://127.0.0.1:3000
```

### 3. Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/kids` | Retrieve all kids |
| GET | `/kids/{id}` | Retrieve kid by ID |
| POST | `/kids` | Create new kid |
| PUT | `/kids/{id}` | Update existing kid |
| DELETE | `/kids/{id}` | Delete kid |

## 🧪 Testing

### cURL
```bash
# Get all kids
curl -X GET http://127.0.0.1:3000/kids

# Get kid by ID=1
curl -X GET http://127.0.0.1:3000/kids/1

# Create new kid
curl -X POST http://127.0.0.1:3000/kids \
  -H "Content-Type: application/json" \
  -d '{"name": "John Smith", "birthdate": "2015-03-15"}'

# Update kid
curl -X PUT http://127.0.0.1:3000/kids/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Smith Updated", "birthdate": "2015-03-15"}'

# Delete kid
curl -X DELETE http://127.0.0.1:3000/kids/1
```

### Postman
Import collections from the `postman/` folder into Postman:
- `postman/kid_service.json` - Kid Service CRUD operations
- `postman/caregiver_service.json` - Caregiver Service CRUD + validation endpoints

Both collections contain all endpoints configured for local environment.

## 🔧 Development Workflow

### Code Modifications
1. Edit Go files in `cmd/kid-service/` or `internal/kidhandler/`
2. Rebuild: `mage build:kid`
3. Changes are automatically detected by SAM (no restart needed)

### Debugging
```bash
# SAM local with debug logs
sam local start-api --port 3000 --debug

# Check Lambda logs
# Logs appear in the terminal where SAM is running
```

## 📁 Project Structure

```
astras-mono-api/
├── cmd/kid-service/          # Main Lambda handler
├── internal/kidhandler/      # Service business logic
├── bin/kid-service/          # Compiled binaries
├── template.yaml             # SAM configuration
├── postman_collection.json   # Postman collection
└── LOCAL_DEVELOPMENT.md      # This document
```

## 🎯 Key Files

- **`template.yaml`** - AWS SAM configuration with API Gateway + Lambda definition
- **`cmd/kid-service/main.go`** - Lambda function entry point
- **`internal/kidhandler/handler.go`** - CRUD operations implementation
- **`bin/kid-service/bootstrap`** - Compiled Lambda binary

## 🚨 Troubleshooting

### SAM won't start
```bash
# Check if Docker is running
docker --version

# Check ports
lsof -i :3000
```

### Build errors
```bash
# Clear cache and rebuild
go clean -cache
mage clean:build
mage build:kid
```

### Go import errors
```bash
# Check Go modules
go mod tidy
go mod verify
```

## 💡 Tips

1. **Code changes** - After modifying code, run `mage build:kid`, SAM will automatically detect new binary
2. **CORS** - Configured in `template.yaml` for all origins
3. **Ports** - Default port 3000, can be changed with `--port`
4. **Performance** - SAM local may be slower than production due to Docker overhead
5. **Hot Reload** - SAM local doesn't support hot reload, need to rebuild binary

## 🔗 Useful Links

- [AWS SAM CLI Documentation](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
- [SAM Local Testing](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-cli-command-reference-sam-local-start-api.html)
- [Go Lambda Runtime](https://github.com/aws/aws-lambda-go)