#!/bin/bash

# Database migration management script for Astras API
# Usage: ./scripts/migrate.sh [command] [stage] [options]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DEFAULT_STAGE="dev"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    cat << EOF
Astras Database Migration Manager

USAGE:
    $0 <command> [stage] [options]

COMMANDS:
    deploy-infra    Deploy RDS infrastructure (required first step)
    deploy-service  Deploy migration Lambda service
    migrate         Run all pending migrations
    rollback        Rollback last migration(s)
    status          Show current migration status
    create          Create a new migration file
    local-migrate   Run migrations against local database
    local-status    Show migration status for local database

STAGES:
    dev (default)   Development environment
    staging         Staging environment  
    prod            Production environment

OPTIONS:
    --steps N       Number of rollback steps (default: 1)
    --name NAME     Migration name for create command
    --help, -h      Show this help message

EXAMPLES:
    $0 deploy-infra dev
    $0 migrate dev
    $0 rollback dev --steps 2
    $0 create --name "add_user_table"
    $0 local-migrate

WORKFLOW:
    1. First time setup:
       $0 deploy-infra dev
       $0 deploy-service dev
       $0 migrate dev

    2. Regular migrations:
       $0 create --name "your_migration_name"
       # Edit the generated files in database/migrations/
       $0 deploy-service dev  # Update Lambda with new migrations
       $0 migrate dev

    3. Rollbacks:
       $0 rollback dev --steps 1

EOF
}

# Parse command line arguments
COMMAND=$1
STAGE=${2:-$DEFAULT_STAGE}

# Parse options
while [[ $# -gt 0 ]]; do
    case $1 in
        --steps)
            STEPS="$2"
            shift 2
            ;;
        --name)
            MIGRATION_NAME="$2"
            shift 2
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            shift
            ;;
    esac
done

# Validate stage
if [[ ! "$STAGE" =~ ^(dev|staging|prod)$ ]]; then
    log_error "Invalid stage: $STAGE. Must be one of: dev, staging, prod"
    exit 1
fi

# Set environment variables
export STAGE="$STAGE"
export AWS_PROFILE="${AWS_PROFILE:-default}"

case $COMMAND in
    deploy-infra)
        log_info "Deploying RDS infrastructure for stage: $STAGE"
        cd "$PROJECT_ROOT"
        
        if ! command -v serverless &> /dev/null; then
            log_error "Serverless Framework not installed. Install with: npm install -g serverless"
            exit 1
        fi
        
        # Check if DB_PASSWORD is set
        if [[ -z "$DB_PASSWORD" ]]; then
            log_warning "DB_PASSWORD not set. Using default password."
            log_warning "For production, set DB_PASSWORD environment variable!"
            export DB_PASSWORD="DefaultPassword123!"
        fi
        
        serverless deploy --config serverless-infrastructure.yml --stage "$STAGE" --verbose
        log_success "Infrastructure deployed successfully for stage: $STAGE"
        ;;
        
    deploy-service)
        log_info "Building and deploying migration service for stage: $STAGE"
        cd "$PROJECT_ROOT"
        
        # Build migration service
        if ! command -v mage &> /dev/null; then
            log_info "Installing Mage..."
            go install github.com/magefile/mage@latest
        fi
        
        mage build:migration
        
        # Deploy service
        cd "services/migration-service"
        serverless deploy --stage "$STAGE" --verbose
        cd "$PROJECT_ROOT"
        
        log_success "Migration service deployed successfully for stage: $STAGE"
        ;;
        
    migrate)
        log_info "Running database migrations for stage: $STAGE"
        
        # Get the API endpoint from CloudFormation outputs
        API_ID=$(aws cloudformation describe-stacks \
            --stack-name "astras-migration-service-$STAGE" \
            --query 'Stacks[0].Outputs[?OutputKey==`HttpApiId`].OutputValue' \
            --output text 2>/dev/null || echo "")
            
        if [[ -z "$API_ID" ]]; then
            log_error "Migration service not found. Deploy it first with: $0 deploy-service $STAGE"
            exit 1
        fi
        
        REGION=$(aws configure get region || echo "eu-central-1")
        API_URL="https://$API_ID.execute-api.$REGION.amazonaws.com/migrations/migrate"
        
        log_info "Calling migration API: $API_URL"
        
        RESPONSE=$(curl -s -X POST \
            -H "Content-Type: application/json" \
            "$API_URL" || echo '{"success": false, "error": "API call failed"}')
            
        echo "$RESPONSE" | jq '.'
        
        if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
            log_success "Migrations completed successfully"
        else
            log_error "Migration failed"
            exit 1
        fi
        ;;
        
    rollback)
        STEPS=${STEPS:-1}
        log_info "Rolling back $STEPS migration(s) for stage: $STAGE"
        
        # Get the API endpoint from CloudFormation outputs
        API_ID=$(aws cloudformation describe-stacks \
            --stack-name "astras-migration-service-$STAGE" \
            --query 'Stacks[0].Outputs[?OutputKey==`HttpApiId`].OutputValue' \
            --output text 2>/dev/null || echo "")
            
        if [[ -z "$API_ID" ]]; then
            log_error "Migration service not found. Deploy it first with: $0 deploy-service $STAGE"
            exit 1
        fi
        
        REGION=$(aws configure get region || echo "eu-central-1")
        API_URL="https://$API_ID.execute-api.$REGION.amazonaws.com/migrations/rollback"
        
        log_info "Calling rollback API: $API_URL (steps: $STEPS)"
        
        RESPONSE=$(curl -s -X POST \
            -H "Content-Type: application/json" \
            -d "{\"steps\": $STEPS}" \
            "$API_URL" || echo '{"success": false, "error": "API call failed"}')
            
        echo "$RESPONSE" | jq '.'
        
        if echo "$RESPONSE" | jq -e '.success' > /dev/null; then
            log_success "Rollback completed successfully"
        else
            log_error "Rollback failed"
            exit 1
        fi
        ;;
        
    status)
        log_info "Checking migration status for stage: $STAGE"
        
        # Get the API endpoint from CloudFormation outputs
        API_ID=$(aws cloudformation describe-stacks \
            --stack-name "astras-migration-service-$STAGE" \
            --query 'Stacks[0].Outputs[?OutputKey==`HttpApiId`].OutputValue' \
            --output text 2>/dev/null || echo "")
            
        if [[ -z "$API_ID" ]]; then
            log_error "Migration service not found. Deploy it first with: $0 deploy-service $STAGE"
            exit 1
        fi
        
        REGION=$(aws configure get region || echo "eu-central-1")
        API_URL="https://$API_ID.execute-api.$REGION.amazonaws.com/migrations/status"
        
        log_info "Calling status API: $API_URL"
        
        RESPONSE=$(curl -s -X GET \
            -H "Content-Type: application/json" \
            "$API_URL" || echo '{"success": false, "error": "API call failed"}')
            
        echo "$RESPONSE" | jq '.'
        ;;
        
    create)
        if [[ -z "$MIGRATION_NAME" ]]; then
            log_error "Migration name is required. Use --name option."
            exit 1
        fi
        
        log_info "Creating new migration: $MIGRATION_NAME"
        
        # Generate timestamp-based version
        VERSION=$(date +"%Y%m%d%H%M%S")
        MIGRATIONS_DIR="$PROJECT_ROOT/database/migrations"
        
        # Create up migration
        UP_FILE="$MIGRATIONS_DIR/${VERSION}_${MIGRATION_NAME}.up.sql"
        cat > "$UP_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created: $(date)
-- Description: [Add description here]

-- Add your up migration SQL here
-- Example:
-- CREATE TABLE example_table (
--     id SERIAL PRIMARY KEY,
--     name VARCHAR(255) NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

EOF

        # Create down migration
        DOWN_FILE="$MIGRATIONS_DIR/${VERSION}_${MIGRATION_NAME}.down.sql"
        cat > "$DOWN_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created: $(date)
-- Description: Rollback for $MIGRATION_NAME

-- Add your down migration SQL here
-- Example:
-- DROP TABLE IF EXISTS example_table;

EOF

        log_success "Created migration files:"
        echo "  Up:   $UP_FILE"
        echo "  Down: $DOWN_FILE"
        log_info "Edit these files with your SQL, then deploy with: $0 deploy-service $STAGE"
        ;;
        
    local-migrate)
        log_info "Running migrations against local database"
        cd "$PROJECT_ROOT"
        
        # Build and run local migration tool
        if [[ ! -f "bin/migration-service/bootstrap" ]]; then
            log_info "Building migration service locally..."
            mage build:migrationLocal
        fi
        
        # Set local database environment variables
        export DB_HOST="localhost"
        export DB_PORT="5432"
        export DB_NAME="astras"
        export DB_USER="postgres"
        export DB_PASSWORD="password"
        export DB_SSL_MODE="disable"
        
        # Check if local database is running
        if ! nc -z localhost 5432 2>/dev/null; then
            log_error "Local PostgreSQL not running. Start with: docker-compose up -d"
            exit 1
        fi
        
        # Run migrations locally
        ./bin/migration-service/bootstrap local-migrate
        log_success "Local migrations completed"
        ;;
        
    local-status)
        log_info "Checking local database migration status"
        cd "$PROJECT_ROOT"
        
        # Set local database environment variables
        export DB_HOST="localhost" 
        export DB_PORT="5432"
        export DB_NAME="astras"
        export DB_USER="postgres"
        export DB_PASSWORD="password"
        export DB_SSL_MODE="disable"
        
        # Check if local database is running
        if ! nc -z localhost 5432 2>/dev/null; then
            log_error "Local PostgreSQL not running. Start with: docker-compose up -d"
            exit 1
        fi
        
        # Build if needed
        if [[ ! -f "bin/migration-service/bootstrap" ]]; then
            log_info "Building migration service locally..."
            mage build:migrationLocal
        fi
        
        ./bin/migration-service/bootstrap local-status
        ;;
        
    *)
        log_error "Unknown command: $COMMAND"
        echo ""
        show_help
        exit 1
        ;;
esac