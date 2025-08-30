# Claude Instructions for astras-mono-api

## Project Overview
This is the astras-mono-api project - a Go monorepo for API services deployed to AWS using the Serverless Framework.

## Development Guidelines

### Code Style
- Follow Go idioms and best practices
- Use `gofmt` for consistent formatting
- Follow the Go Code Review Comments guidelines
- Write clean, readable, and maintainable code
- Use meaningful variable and function names

### Testing
- Write tests for new features and bug fixes
- Ensure all tests pass before marking tasks as complete
- Run test suite with: `go test ./...`
- Aim for good test coverage

### Linting & Formatting
- Run formatting: `go fmt ./...`
- Run linting: `golangci-lint run` (if configured)
- Run vet: `go vet ./...`
- Fix all linting errors before completing tasks

### Git Workflow
- Create descriptive commit messages
- Keep commits focused and atomic
- Don't commit directly unless explicitly asked

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

### Deployment
- Use Serverless Framework for AWS Lambda deployment
- Each service has its own `serverless.yml` configuration
- Build binaries for Linux before deployment
- Use AWS API Gateway for HTTP endpoints

## Notes
- This is a monorepo structure - consider impacts across services
- Each service should be in its own directory under cmd/
- Shared packages should be in internal/ or pkg/
- [Add project-specific notes and reminders here]