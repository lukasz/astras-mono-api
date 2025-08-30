# Claude Instructions for astras-mono-api

## Communication Guidelines
- **Always respond in the same language the user writes in** - If user writes in Polish, respond in Polish. If user writes in English, respond in English.
- **Documentation and code comments should always be in English** regardless of conversation language
- Keep responses concise and direct

## Project Overview
This is the astras-mono-api project - a Go monorepo for API services deployed to AWS using the Serverless Framework.

## Development Guidelines

### Code Style
- Follow Go idioms and best practices
- Use `gofmt` for consistent formatting
- Follow the Go Code Review Comments guidelines
- Write clean, readable, and maintainable code
- Use meaningful variable and function names
- **Define constants for limits and formats**: Use named constants instead of magic numbers for validation limits, format strings, and configuration values
- **Constant naming**: Follow Go convention with proper comments (e.g., `// MaxKidAge defines the maximum allowed age...`)

### Testing
- Write tests for new features and bug fixes
- Ensure all tests pass before marking tasks as complete
- Run test suite with: `go test ./...`
- Aim for good test coverage
- **Use JSON fixtures for test data**: Store test cases in JSON files within `testdata/fixtures/` directories for maintainability and reusability
- **Deterministic testing**: Use time injection parameters instead of `time.Now()` for predictable test behavior
- **Test data organization**: Group related test cases in fixture files (e.g., `model_validation_tests.json`, `model_calculation_tests.json`)

### Linting & Formatting
- Run formatting: `go fmt ./...`
- Run linting: `golangci-lint run` (if configured)
- Run vet: `go vet ./...`
- Fix all linting errors before completing tasks

### Git Workflow
- Create descriptive commit messages
- Keep commits focused and atomic
- Don't commit directly unless explicitly asked
- **Split commits logically**: When committing multiple changes, create separate atomic commits for each logical change based on conversation history
- **Do NOT add AI-generated attribution** to commit messages unless specifically requested
- **ALWAYS check before committing** that commonly recognized files/folders/patterns that should NOT be committed are properly excluded:
  - `node_modules/`, `vendor/` (dependency directories)
  - `.env`, `.env.*` (environment/secret files)
  - Build artifacts: `bin/`, `dist/`, `build/`, `target/`
  - IDE files: `.vscode/`, `.idea/`, `*.swp`
  - OS files: `.DS_Store`, `Thumbs.db`
  - Logs: `*.log`, `logs/`
  - Temporary files: `*.tmp`, `*.temp`
  - Database files: `*.db`, `*.sqlite`
  - Compiled binaries and executables
  - Large media files that should use Git LFS
  - API keys, certificates, or any sensitive data

### Project Structure
- [Add project structure details here as the project develops]
- Follow standard Go project layout for monorepo

### API Conventions
- Use RESTful API design patterns
- Implement proper error handling
- Use structured logging
- [Add specific API conventions here]

### Dependencies
- Check go.mod before adding new dependencies
- Use Go modules for dependency management
- Prefer standard library when possible
- Document any new dependencies added

### Environment Variables
- [Document required environment variables here]
- Use `.env` files for local development

### Build System
- Use Mage for building and deployment automation
- Install Mage: `go install github.com/magefile/mage@latest`
- Use `mage -l` to list available targets
- Build targets:
  - `mage build:all` - Build all services
  - `mage build:kid` - Build kid service
  - `mage build:caregiver` - Build caregiver service
  - `mage build:star` - Build star service

### Local Development
- **Use AWS SAM CLI for local testing** - Provides accurate simulation of AWS serverless architecture
- Install SAM CLI: `brew tap aws/tap && brew install aws-sam-cli`
- Local development workflow:
  1. Build service: `export PATH=$PATH:$(go env GOPATH)/bin && mage build:kid`
  2. Start SAM local: `sam local start-api --port 3000`
  3. Test endpoints: http://127.0.0.1:3000/kids
  4. Use Postman collection: `postman_collection.json`
- **Key files for local dev**:
  - `template.yaml` - SAM template with API Gateway + Lambda configuration
  - `LOCAL_DEVELOPMENT.md` - Detailed local development guide
  - `postman_collection.json` - Ready-to-import Postman collection

### Deployment
- Use Serverless Framework for AWS Lambda deployment
- Each service has its own `serverless.yml` configuration
- Build binaries for Linux before deployment using Mage
- Use AWS API Gateway for HTTP endpoints
- **Local testing with SAM before deployment is recommended**
- Deploy targets:
  - `mage deploy:all` - Deploy all services
  - `mage deploy:kid` - Deploy kid service
  - `mage deploy:caregiver` - Deploy caregiver service
  - `mage deploy:star` - Deploy star service

### Common Mage Commands
```bash
# List available targets
mage -l

# Build all services
mage build:all

# Build specific service
mage build:kid
mage build:caregiver
mage build:star

# Deploy all services
mage deploy:all

# Deploy specific service
mage deploy:kid
mage deploy:caregiver
mage deploy:star

# Test
mage test:all
mage test:coverage

# Clean
mage clean:all
mage clean:build
mage clean:deploy

# Code quality
mage format
mage lint
mage tidy

# List services
mage services
```

### Local SAM Commands
```bash
# Start local API Gateway + Lambda
sam local start-api --port 3000

# Start local API with debug logs
sam local start-api --port 3000 --debug

# Invoke specific function directly
sam local invoke KidFunction --event test-event.json

# Generate sample events for testing
sam local generate-event apigateway aws-proxy

# Validate SAM template
sam validate
```

### npm Scripts (Alternative)
```bash
# Alternative using npm scripts
npm run build          # Build all services
npm run build:kid      # Build kid service
npm run deploy         # Deploy all services
npm run deploy:kid     # Deploy kid service
npm test              # Run tests
npm run clean         # Clean artifacts
```

## Project Structure
```
astras-mono-api/
├── cmd/
│   ├── kid-service/     # Kid service Lambda handler
│   ├── caregiver-service/ # Caregiver service Lambda handler
│   └── star-service/    # Star service Lambda handler
├── services/
│   ├── kid-service/     # Kid service Serverless config
│   ├── caregiver-service/ # Caregiver service Serverless config
│   └── star-service/    # Star service Serverless config
├── internal/            # Private application code
├── pkg/                 # Public library code
├── bin/                 # Built binaries (ignored by git)
├── template.yaml        # AWS SAM template for local development
├── postman_collection.json # Postman collection for API testing
├── LOCAL_DEVELOPMENT.md # Local development guide
├── magefile.go         # Mage build configuration
└── package.json        # Node.js dependencies for Serverless
```

## Services
- **kid-service**: Manages children/kids in the system
- **caregiver-service**: Manages caregivers and guardians
- **star-service**: Manages star rewards and achievements

## Notes
- This is a monorepo structure - consider impacts across services
- Each service is built as a standalone Lambda function
- Services communicate via API calls when needed
- Shared code should be placed in internal/ or pkg/ directories
- **Local development uses AWS SAM CLI** for accurate serverless simulation
- Use `template.yaml` for SAM configuration, `serverless.yml` for production deployment
- See `LOCAL_DEVELOPMENT.md` for detailed local development instructions