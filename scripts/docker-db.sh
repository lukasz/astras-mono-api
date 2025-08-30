#!/bin/bash

# Docker database management script for Astras API

set -e

COMPOSE_FILE="docker-compose.yml"

case "$1" in
    "start")
        echo "🚀 Starting PostgreSQL database..."
        docker-compose up -d postgres
        echo "✅ PostgreSQL is starting up. Waiting for health check..."
        docker-compose exec postgres pg_isready -U postgres -d astras
        echo "🎉 Database is ready!"
        echo ""
        echo "📊 Connection details:"
        echo "  Host: localhost"
        echo "  Port: 5432"
        echo "  Database: astras"
        echo "  Username: postgres"
        echo "  Password: password"
        ;;
    "stop")
        echo "🛑 Stopping PostgreSQL database..."
        docker-compose stop postgres
        echo "✅ Database stopped"
        ;;
    "restart")
        echo "🔄 Restarting PostgreSQL database..."
        docker-compose restart postgres
        echo "✅ Database restarted"
        ;;
    "logs")
        echo "📝 Showing PostgreSQL logs..."
        docker-compose logs -f postgres
        ;;
    "shell")
        echo "🐚 Opening PostgreSQL shell..."
        docker-compose exec postgres psql -U postgres -d astras
        ;;
    "admin")
        echo "🚀 Starting PostgreSQL with pgAdmin..."
        docker-compose --profile admin up -d
        echo "✅ Database and pgAdmin started!"
        echo ""
        echo "📊 Access points:"
        echo "  Database: localhost:5432"
        echo "  pgAdmin: http://localhost:8080"
        echo "    Email: admin@astras.local"
        echo "    Password: admin"
        ;;
    "clean")
        echo "🧹 Cleaning up database (this will delete all data)..."
        read -p "Are you sure? This will delete all data. [y/N]: " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose down -v
            docker volume rm astras-mono-api_postgres_data astras-mono-api_pgadmin_data 2>/dev/null || true
            echo "✅ Database cleaned"
        else
            echo "❌ Cancelled"
        fi
        ;;
    "status")
        echo "📊 Database status:"
        docker-compose ps postgres
        ;;
    *)
        echo "🗄️  Astras API Database Management"
        echo ""
        echo "Usage: $0 {start|stop|restart|logs|shell|admin|clean|status}"
        echo ""
        echo "Commands:"
        echo "  start   - Start PostgreSQL database"
        echo "  stop    - Stop PostgreSQL database"
        echo "  restart - Restart PostgreSQL database"
        echo "  logs    - Show database logs"
        echo "  shell   - Open PostgreSQL shell (psql)"
        echo "  admin   - Start database with pgAdmin web interface"
        echo "  clean   - Remove database and all data (destructive)"
        echo "  status  - Show database container status"
        echo ""
        echo "Examples:"
        echo "  $0 start       # Start database"
        echo "  $0 shell       # Connect to database"
        echo "  $0 admin       # Start with web admin"
        ;;
esac