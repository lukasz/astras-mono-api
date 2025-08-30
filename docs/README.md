# Astras Mono API

Go monorepo for Astras API services deployed on AWS Lambda.

## 🚀 Quick Start

### Local Development with Database
```bash
# 1. Start PostgreSQL database
docker-compose up -d

# 2. Build service
export PATH=$PATH:$(go env GOPATH)/bin
mage build:kid

# 3. Start locally with AWS SAM (with database connection)
sam local start-api --env-vars env.json --docker-network astras-mono-api_astras-network --port 3000

# 4. Test API
curl http://127.0.0.1:3000/kids
```

### Documentation
- **[LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md)** - Complete local development guide
- **[DATABASE.md](DATABASE.md)** - Database architecture and setup guide
- **[CLAUDE.md](CLAUDE.md)** - Instructions for AI assistant
- **[postman/](postman/)** - Postman collections for API testing

## 📦 Services
- **kid-service** - Manages children/kids in the system ✅ *Database integrated*
- **caregiver-service** - Manages caregivers and guardians  
- **star-service** - Manages star rewards and achievements

## 🏗️ Architecture

### Database Layer
- **PostgreSQL** as primary database with pgx/v5 driver
- **Repository pattern** for clean data access layer
- **Database migrations** in `database/migrations/`
- **Connection pooling** and health checks
- **Environment-based configuration**

### Models
- **Kids** - Children with birthdate and age calculation
- **Caregivers** - Adults responsible for kids
- **Transactions** - Star earning/spending records

### Local Development
- **Docker Compose** for PostgreSQL with sample data
- **Environment variables** via `env.json` for SAM local
- **Database schema** auto-created on container startup

## 🛠️ Tech Stack
- **Go 1.23+** - Backend language
- **PostgreSQL** - Primary database (AWS RDS in production)
- **pgx/v5** - PostgreSQL driver with sqlx
- **Docker & Docker Compose** - Local database development
- **AWS Lambda** - Serverless compute
- **AWS SAM CLI** - Local development and testing
- **Mage** - Build automation
- **Serverless Framework** - Deployment
