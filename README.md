# Astras Mono API

Go monorepo for Astras API services deployed on AWS Lambda.

## 🚀 Quick Start

### Local Development
```bash
# 1. Build service
export PATH=$PATH:$(go env GOPATH)/bin
mage build:kid

# 2. Start locally with AWS SAM
sam local start-api --port 3000

# 3. Test API
curl http://127.0.0.1:3000/kids
```

### Documentation
- **[LOCAL_DEVELOPMENT.md](LOCAL_DEVELOPMENT.md)** - Complete local development guide
- **[CLAUDE.md](CLAUDE.md)** - Instructions for AI assistant
- **[postman_collection.json](postman_collection.json)** - Postman collection for testing

## 📦 Services
- **kid-service** - Manages children/kids in the system
- **caregiver-service** - Manages caregivers and guardians  
- **star-service** - Manages star rewards and achievements

## 🛠️ Tech Stack
- **Go 1.21+** - Backend language
- **AWS Lambda** - Serverless compute
- **AWS SAM CLI** - Local development
- **Mage** - Build automation
- **Serverless Framework** - Deployment
