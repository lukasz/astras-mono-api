#!/bin/bash

# Local logging setup for Astras API development
# This script sets up local logging infrastructure using Docker

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
    $0 <command> [options]

COMMANDS:
    start       Start local logging infrastructure
    stop        Stop all logging services
    logs        View aggregated logs
    dashboard   Open local dashboard (if available)
    clean       Clean up log data
    status      Check status of logging services

OPTIONS:
    --elastic   Use ELK stack instead of simple file logging
    --help, -h  Show this help message

EXAMPLES:
    $0 start                    # Start simple file-based logging
    $0 start --elastic         # Start with Elasticsearch + Kibana
    $0 logs                    # View recent logs
    $0 stop                    # Stop all services

EOF
}

# Parse command line arguments
COMMAND=$1
USE_ELASTIC=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --elastic)
            USE_ELASTIC=true
            shift
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

case $COMMAND in
    start)
        log_info "Starting local logging infrastructure..."
        
        if [ "$USE_ELASTIC" = true ]; then
            log_info "Using ELK Stack (Elasticsearch + Kibana + Logstash)"
            
            # Create ELK docker-compose file
            cat > "$PROJECT_ROOT/docker-compose.logging.yml" << 'EOF'
version: '3.8'

services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    container_name: astras-elasticsearch
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"
      - xpack.security.enabled=false
      - xpack.security.enrollment.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data
    networks:
      - astras-logging

  kibana:
    image: docker.elastic.co/kibana/kibana:8.11.0
    container_name: astras-kibana
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
      - xpack.security.enabled=false
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    networks:
      - astras-logging

  logstash:
    image: docker.elastic.co/logstash/logstash:8.11.0
    container_name: astras-logstash
    volumes:
      - ./logstash/pipeline:/usr/share/logstash/pipeline
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml
      - ./logs:/logs
    ports:
      - "5044:5044"
      - "9600:9600"
    environment:
      - "LS_JAVA_OPTS=-Xmx256m -Xms256m"
    depends_on:
      - elasticsearch
    networks:
      - astras-logging

volumes:
  elasticsearch_data:

networks:
  astras-logging:
    external: true
EOF

            # Create Logstash configuration
            mkdir -p "$PROJECT_ROOT/logstash/pipeline"
            mkdir -p "$PROJECT_ROOT/logstash/config"
            
            cat > "$PROJECT_ROOT/logstash/config/logstash.yml" << 'EOF'
http.host: "0.0.0.0"
xpack.monitoring.elasticsearch.hosts: [ "http://elasticsearch:9200" ]
EOF

            cat > "$PROJECT_ROOT/logstash/pipeline/logstash.conf" << 'EOF'
input {
  file {
    path => "/logs/astras-*.log"
    start_position => "beginning"
    codec => "json"
    tags => ["astras"]
  }
}

filter {
  if [tags] and "astras" in [tags] {
    # Parse timestamp
    date {
      match => [ "@timestamp", "ISO8601" ]
    }
    
    # Add parsed fields
    if [level] {
      mutate { add_tag => [ "level_%{level}" ] }
    }
    
    if [service] {
      mutate { add_tag => [ "service_%{service}" ] }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "astras-logs-%{+YYYY.MM.dd}"
  }
  
  stdout { 
    codec => rubydebug 
  }
}
EOF

            # Create logging network
            docker network create astras-logging 2>/dev/null || true
            
            # Start ELK stack
            cd "$PROJECT_ROOT"
            docker-compose -f docker-compose.logging.yml up -d
            
            log_success "ELK Stack started!"
            log_info "Elasticsearch: http://localhost:9200"
            log_info "Kibana: http://localhost:5601"
            
        else
            log_info "Using simple file-based logging"
            
            # Create logs directory
            mkdir -p "$PROJECT_ROOT/logs"
            
            # Create simple docker-compose for log aggregation
            cat > "$PROJECT_ROOT/docker-compose.logging.yml" << 'EOF'
version: '3.8'

services:
  log-aggregator:
    image: fluent/fluent-bit:latest
    container_name: astras-log-aggregator
    volumes:
      - ./logs:/logs
      - ./fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf
    ports:
      - "24224:24224"
    networks:
      - astras-logging

  log-viewer:
    image: goharbor/harbor-log:latest
    container_name: astras-log-viewer
    volumes:
      - ./logs:/logs:ro
    ports:
      - "8080:8080"
    networks:
      - astras-logging

networks:
  astras-logging:
    external: true

EOF

            # Create Fluent Bit configuration
            cat > "$PROJECT_ROOT/fluent-bit.conf" << 'EOF'
[SERVICE]
    Flush         1
    Log_Level     info
    Daemon        off
    Parsers_File  parsers.conf

[INPUT]
    Name              tail
    Path              /logs/astras-*.log
    Parser            json
    Tag               astras.*
    Refresh_Interval  5

[OUTPUT]
    Name  file
    Match *
    Path  /logs/
    File  aggregated.log
    Format json_lines

[OUTPUT]
    Name  stdout
    Match *
EOF

            # Create logging network
            docker network create astras-logging 2>/dev/null || true
            
            cd "$PROJECT_ROOT"
            docker-compose -f docker-compose.logging.yml up -d
            
            log_success "File-based logging started!"
            log_info "Logs directory: $PROJECT_ROOT/logs"
            log_info "Log viewer: http://localhost:8080"
        fi
        ;;
        
    stop)
        log_info "Stopping local logging infrastructure..."
        cd "$PROJECT_ROOT"
        
        if [ -f "docker-compose.logging.yml" ]; then
            docker-compose -f docker-compose.logging.yml down
            log_success "Logging services stopped"
        else
            log_warning "No logging services found to stop"
        fi
        ;;
        
    logs)
        log_info "Viewing recent logs..."
        
        if [ -f "$PROJECT_ROOT/logs/aggregated.log" ]; then
            tail -f "$PROJECT_ROOT/logs/aggregated.log" | jq -r '"\(.["@timestamp"]) [\(.level)] \(.service): \(.message)"'
        elif [ -d "$PROJECT_ROOT/logs" ]; then
            find "$PROJECT_ROOT/logs" -name "*.log" -exec tail -f {} +
        else
            log_error "No logs found. Start logging infrastructure first."
        fi
        ;;
        
    dashboard)
        log_info "Opening local dashboard..."
        
        if [ "$USE_ELASTIC" = true ]; then
            if command -v open >/dev/null 2>&1; then
                open "http://localhost:5601"
            else
                log_info "Open http://localhost:5601 in your browser"
            fi
        else
            if command -v open >/dev/null 2>&1; then
                open "http://localhost:8080"
            else
                log_info "Open http://localhost:8080 in your browser"
            fi
        fi
        ;;
        
    clean)
        log_warning "This will delete all local log data. Continue? (y/N)"
        read -r response
        if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
            rm -rf "$PROJECT_ROOT/logs"
            docker volume rm $(docker volume ls -q | grep -E "(elasticsearch_data|astras.*log)") 2>/dev/null || true
            log_success "Local log data cleaned"
        else
            log_info "Cancelled"
        fi
        ;;
        
    status)
        log_info "Checking logging services status..."
        
        cd "$PROJECT_ROOT"
        if [ -f "docker-compose.logging.yml" ]; then
            docker-compose -f docker-compose.logging.yml ps
        else
            log_warning "No logging infrastructure configured"
        fi
        ;;
        
    *)
        log_error "Unknown command: $COMMAND"
        echo ""
        show_help
        exit 1
        ;;
esac