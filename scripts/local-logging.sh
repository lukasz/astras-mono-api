#!/bin/bash

# Local logging setup for Astras API development
# Simple Docker-based log viewer using Dozzle

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

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
Local Logging Setup for Astras API

USAGE:
    $0 <command>

COMMANDS:
    start       Start Dozzle log viewer
    stop        Stop logging services
    logs        View logs in terminal
    dashboard   Open Dozzle web interface
    clean       Clean up log data
    status      Check status of services

EXAMPLES:
    $0 start                    # Start Dozzle log viewer
    $0 logs                    # View recent logs in terminal  
    $0 dashboard               # Open Dozzle at http://localhost:8080
    $0 stop                    # Stop all services

EOF
}

# Parse command line arguments  
COMMAND=$1

if [[ "$COMMAND" == "--help" || "$COMMAND" == "-h" ]]; then
    show_help
    exit 0
fi

case $COMMAND in
    start)
        log_info "Starting Dozzle log viewer..."
        
        # Create logs directory
        mkdir -p "$PROJECT_ROOT/logs"
        
        # Create logging network
        docker network create astras-logging 2>/dev/null || true
        
        cd "$PROJECT_ROOT"
        docker-compose -f config/docker-compose.dozzle.yml up -d
        
        log_success "Dozzle started!"
        log_info "Dozzle: http://localhost:8080 (Docker container logs)"
        log_info "Log files: http://localhost:8081 (File-based logs)"
        ;;
        
    stop)
        log_info "Stopping Dozzle log viewer..."
        cd "$PROJECT_ROOT"
        
        if [ -f "config/docker-compose.dozzle.yml" ]; then
            docker-compose -f config/docker-compose.dozzle.yml down  
            log_success "Dozzle stopped"
        else
            log_warning "No Dozzle services found to stop"
        fi
        ;;
        
    logs)
        log_info "Viewing recent logs in terminal..."
        
        if [ -d "$PROJECT_ROOT/logs" ]; then
            find "$PROJECT_ROOT/logs" -name "*.log" -exec tail -f {} +
        else
            log_error "No logs found. Logs directory: $PROJECT_ROOT/logs"
        fi
        ;;
        
    dashboard)
        log_info "Opening Dozzle dashboard..."
        
        if command -v open >/dev/null 2>&1; then
            open "http://localhost:8080"
        else
            log_info "Open http://localhost:8080 in your browser"
        fi
        ;;
        
    clean)
        log_warning "This will delete all local log data. Continue? (y/N)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            rm -rf "$PROJECT_ROOT/logs"
            log_success "Local log data cleaned"
        else
            log_info "Cancelled"
        fi
        ;;
        
    status)
        log_info "Checking Dozzle service status..."
        cd "$PROJECT_ROOT"
        
        if [ -f "config/docker-compose.dozzle.yml" ]; then
            docker-compose -f config/docker-compose.dozzle.yml ps
        else
            log_warning "No Dozzle configuration found"
        fi
        ;;
        
    *)
        log_error "Unknown command: $COMMAND"
        echo ""
        show_help
        exit 1
        ;;
esac