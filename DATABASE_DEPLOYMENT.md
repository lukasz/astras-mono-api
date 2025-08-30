# Database Deployment Guide

This guide explains how to deploy and manage the PostgreSQL RDS database infrastructure and run migrations for the Astras API project.

## Overview

The Astras project uses:
- **PostgreSQL RDS** for production database
- **Serverless Framework** for infrastructure deployment
- **Lambda functions** for migration management
- **AWS Systems Manager** for secure configuration storage
- **VPC** for network security

## Prerequisites

1. **AWS CLI** configured with appropriate credentials
2. **Serverless Framework** installed globally: `npm install -g serverless`
3. **Mage** build tool: `go install github.com/magefile/mage@latest`
4. **jq** for JSON processing: `brew install jq` (macOS) or equivalent
5. **curl** for API calls

## Quick Start

### 1. Set Environment Variables

```bash
# Required for production - set a secure password
export DB_PASSWORD="YourSecurePassword123!"

# Optional - AWS profile if not using default
export AWS_PROFILE="your-aws-profile"
```

### 2. Deploy Infrastructure

```bash
# Deploy RDS infrastructure for development
./scripts/migrate.sh deploy-infra dev

# Deploy migration service
./scripts/migrate.sh deploy-service dev

# Run initial migrations
./scripts/migrate.sh migrate dev
```

### 3. Verify Deployment

```bash
# Check migration status
./scripts/migrate.sh status dev
```

## Detailed Deployment Steps

### Step 1: Infrastructure Deployment

The infrastructure includes:
- VPC with public and private subnets
- RDS PostgreSQL instance
- Security groups
- Systems Manager parameters
- Secrets Manager for passwords

```bash
# Deploy to development environment
mage migration:deploy

# Or deploy infrastructure only
mage deploy:infrastructure
```

**Important**: The first deployment takes 10-15 minutes as RDS instance is being created.

### Step 2: Service Deployment

Deploy the migration Lambda service:

```bash
# Build and deploy migration service
mage deploy:migration

# Or manually
mage build:migration
cd services/migration-service
serverless deploy --stage dev
```

### Step 3: Run Migrations

```bash
# Run all pending migrations
./scripts/migrate.sh migrate dev

# Check current status
./scripts/migrate.sh status dev
```

## Migration Management

### Creating New Migrations

```bash
# Create a new migration
./scripts/migrate.sh create --name "add_user_preferences_table"

# This creates two files:
# database/migrations/YYYYMMDDHHMMSS_add_user_preferences_table.up.sql
# database/migrations/YYYYMMDDHHMMSS_add_user_preferences_table.down.sql
```

Edit the generated files with your SQL:

**Up migration** (`*.up.sql`):
```sql
-- Migration: add_user_preferences_table
-- Created: 2024-01-15 10:30:00
-- Description: Add table for user preferences

CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    preference_key VARCHAR(255) NOT NULL,
    preference_value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, preference_key)
);

CREATE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
```

**Down migration** (`*.down.sql`):
```sql
-- Migration: add_user_preferences_table
-- Created: 2024-01-15 10:30:00
-- Description: Rollback for add_user_preferences_table

DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP TABLE IF EXISTS user_preferences;
```

### Deploying New Migrations

After creating migration files:

```bash
# 1. Redeploy migration service to include new files
./scripts/migrate.sh deploy-service dev

# 2. Run migrations
./scripts/migrate.sh migrate dev
```

### Rolling Back Migrations

```bash
# Rollback last migration
./scripts/migrate.sh rollback dev

# Rollback multiple migrations
./scripts/migrate.sh rollback dev --steps 3
```

## Environment-Specific Deployment

### Development Environment

```bash
export DB_PASSWORD="DevPassword123!"
./scripts/migrate.sh deploy-infra dev
./scripts/migrate.sh deploy-service dev
./scripts/migrate.sh migrate dev
```

### Staging Environment

```bash
export DB_PASSWORD="StagingSecurePassword123!"
./scripts/migrate.sh deploy-infra staging
./scripts/migrate.sh deploy-service staging
./scripts/migrate.sh migrate staging
```

### Production Environment

```bash
export DB_PASSWORD="ProductionVerySecurePassword123!"
./scripts/migrate.sh deploy-infra prod
./scripts/migrate.sh deploy-service prod
./scripts/migrate.sh migrate prod
```

## Local Development

For local development, you can run migrations against the local PostgreSQL instance:

```bash
# Start local PostgreSQL
docker-compose up -d

# Run migrations locally
./scripts/migrate.sh local-migrate

# Check local status
./scripts/migrate.sh local-status
```

## Available Mage Commands

```bash
# Infrastructure and migration deployment
mage migration:deploy          # Deploy both infrastructure and migration service
mage deploy:infrastructure     # Deploy only RDS infrastructure
mage deploy:migration         # Deploy only migration service

# Building
mage build:migration          # Build migration service
mage build:migrationLocal     # Build for local development

# General commands
mage migration:status         # Check migration status via API
mage services                 # List all available services
```

## Configuration

### Environment-Specific Settings

The `serverless-infrastructure.yml` file contains environment-specific configurations:

- **Development**: `db.t3.micro`, 20GB storage, no multi-AZ
- **Staging**: `db.t3.small`, 100GB storage, no multi-AZ  
- **Production**: `db.t3.medium`, 200GB storage, multi-AZ enabled

### Security

- Database passwords stored in AWS Systems Manager Parameter Store (encrypted)
- RDS instances deployed in private subnets
- Security groups restrict access to Lambda functions only
- SSL/TLS encryption enabled for database connections

### Backup and Recovery

- **Development**: 7-day backup retention
- **Staging**: 14-day backup retention  
- **Production**: 30-day backup retention, deletion protection enabled

## Troubleshooting

### Common Issues

1. **Infrastructure deployment fails**
   ```bash
   # Check AWS credentials and permissions
   aws sts get-caller-identity
   
   # Verify region setting
   aws configure get region
   ```

2. **Migration API not found**
   ```bash
   # Check if migration service is deployed
   aws cloudformation describe-stacks --stack-name astras-migration-service-dev
   ```

3. **Database connection fails**
   ```bash
   # Verify RDS instance is running
   aws rds describe-db-instances --db-instance-identifier astras-db-dev
   
   # Check security group settings
   aws ec2 describe-security-groups --group-names astras-rds-sg-dev
   ```

4. **Migration files not found**
   ```bash
   # Ensure Lambda layer contains migration files
   # Redeploy migration service
   ./scripts/migrate.sh deploy-service dev
   ```

### Logs and Monitoring

View Lambda function logs:
```bash
# Using AWS CLI
aws logs tail /aws/lambda/astras-migration-service-dev-migration --follow

# Using serverless
cd services/migration-service
serverless logs -f migration --stage dev --tail
```

### Manual Database Access

For debugging, you can access the RDS instance through a bastion host or VPC endpoint. Never expose RDS publicly in production.

## Cost Considerations

- **Development**: ~$15-20/month (db.t3.micro)
- **Staging**: ~$35-45/month (db.t3.small)  
- **Production**: ~$60-80/month (db.t3.medium with multi-AZ)

Additional costs for data transfer and backup storage may apply.

## Best Practices

1. **Always test migrations on development first**
2. **Review migration SQL carefully before deployment**
3. **Use descriptive migration names**
4. **Keep migrations small and focused**
5. **Always write corresponding down migrations**
6. **Monitor migration execution time**
7. **Set strong passwords for production environments**
8. **Use different AWS accounts for different environments**

## Next Steps

After setting up the database infrastructure:

1. **Update service configurations** to use the deployed RDS endpoints
2. **Configure monitoring and alerting** for database metrics
3. **Set up automated backups** verification
4. **Configure log aggregation** for migration activities
5. **Create runbooks** for common operational tasks