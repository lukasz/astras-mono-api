# Database Documentation

## Overview

Astras API uses PostgreSQL as the primary database with pgx/v5 driver for optimal performance and modern PostgreSQL features.

## Architecture

### Technology Stack
- **Database**: PostgreSQL 15+
- **Driver**: pgx/v5/stdlib with sqlx
- **Pattern**: Repository pattern with clean interfaces
- **Migrations**: SQL-based up/down migrations
- **Local Development**: Docker Compose

### Database Schema

#### Tables
1. **kids** - Children/kids in the system
   - `id` (serial, primary key)
   - `name` (varchar(100), not null)
   - `birthdate` (date, not null) - Used to calculate age dynamically
   - `created_at`, `updated_at` (timestamptz)

2. **caregivers** - Adults responsible for kids
   - `id` (serial, primary key)
   - `name` (varchar(100), not null)
   - `email` (varchar(255), unique, not null)
   - `relationship` (enum: parent, guardian, grandparent, relative, caregiver)
   - `created_at`, `updated_at` (timestamptz)

3. **transactions** - Star earning/spending records
   - `id` (serial, primary key)
   - `kid_id` (integer, foreign key to kids)
   - `type` (enum: earn, spend)
   - `amount` (integer, 1-100 stars)
   - `description` (varchar(255), not null)
   - `created_at`, `updated_at` (timestamptz)

## Local Development

### Setup
```bash
# Start PostgreSQL with sample data
docker-compose up -d

# Check database status
docker exec astras-postgres psql -U postgres -d astras -c "\dt"

# View sample data
docker exec astras-postgres psql -U postgres -d astras -c "SELECT name, birthdate, EXTRACT(year FROM age(birthdate)) as age FROM kids;"
```

### Configuration
Database connection is configured via environment variables in `env.json`:

```json
{
  "KidFunction": {
    "DB_HOST": "astras-postgres",
    "DB_PORT": "5432",
    "DB_NAME": "astras",
    "DB_USER": "postgres",
    "DB_PASSWORD": "password",
    "DB_SSL_MODE": "disable"
  }
}
```

### Running with SAM Local
```bash
# Start SAM with database connection
sam local start-api --env-vars env.json --docker-network astras-mono-api_astras-network --port 3000

# Test database integration
curl http://localhost:3000/kids | jq .
```

## Migrations

Located in `database/migrations/` with up/down scripts:

- `001_initial_schema.up.sql` - Creates tables, indexes, triggers
- `001_initial_schema.down.sql` - Drops all objects
- `002_sample_data.up.sql` - Inserts sample data for development
- `002_sample_data.down.sql` - Removes sample data

### Manual Migration
```bash
# Apply schema
docker exec -i astras-postgres psql -U postgres -d astras < database/migrations/001_initial_schema.up.sql

# Add sample data
docker exec -i astras-postgres psql -U postgres -d astras < database/migrations/002_sample_data.up.sql
```

## Repository Pattern

### Interface Layer
All database operations are defined through interfaces in `internal/database/interfaces/`:

```go
type KidRepository interface {
    Create(ctx context.Context, kid *kid.Kid) (*kid.Kid, error)
    GetByID(ctx context.Context, id int) (*kid.Kid, error)
    GetAll(ctx context.Context) ([]*kid.Kid, error)
    Update(ctx context.Context, kid *kid.Kid) (*kid.Kid, error)
    Delete(ctx context.Context, id int) error
    GetByAgeRange(ctx context.Context, minAge, maxAge int) ([]*kid.Kid, error)
}
```

### Implementation Layer
PostgreSQL implementations in `internal/database/postgres/`:
- Connection pooling and health checks
- Structured error handling
- Transaction support
- Query optimization

### Usage in Services
Services inject repository interfaces:

```go
type KidHandler struct {
    repo interfaces.KidRepository
}

func (h *KidHandler) GetAll(ctx context.Context, request events.APIGatewayProxyRequest) (handler.Response, error) {
    kids, err := h.repo.GetAll(ctx)
    // Handle response...
}
```

## Production Considerations

### AWS RDS
- Use PostgreSQL 15+ on AWS RDS for production
- Configure appropriate instance size and storage
- Enable automated backups and point-in-time recovery
- Use Multi-AZ for high availability

### Connection Management
- Connection pooling configured per Lambda function
- Idle connections timeout after 5 minutes
- Maximum 25 concurrent connections per service
- Health checks ensure connection reliability

### Security
- Never commit database passwords to git
- Use AWS Secrets Manager for production credentials
- Enable SSL/TLS for encrypted connections
- Follow principle of least privilege for database users

## Troubleshooting

### Common Issues
1. **Connection refused** - Check if PostgreSQL container is running
2. **Database not found** - Ensure database name matches environment variable
3. **Authentication failed** - Verify username/password in env.json
4. **Column not found** - Ensure migrations have been applied
5. **Network issues** - Verify Docker network configuration

### Debugging
```bash
# Check container status
docker ps | grep postgres

# Check logs
docker logs astras-postgres

# Connect to database directly
docker exec -it astras-postgres psql -U postgres -d astras

# Check SAM logs
tail -f /tmp/sam-*.log
```

## Performance

### Indexes
All tables have appropriate indexes for:
- Primary keys (automatic)
- Foreign keys (transactions.kid_id)
- Frequently queried fields (name, birthdate, created_at)

### Query Optimization
- Use prepared statements through pgx driver
- Connection pooling reduces overhead
- Batch operations where possible
- Monitor slow queries in production

## Data Models

### Kid Model
```go
type Kid struct {
    ID        int       `json:"id" db:"id"`
    Name      string    `json:"name" db:"name"`
    Birthdate time.Time `json:"birthdate" db:"birthdate"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// Age is calculated dynamically from birthdate
func (k *Kid) Age() int { ... }

// MarshalJSON includes computed age field for API compatibility
func (k *Kid) MarshalJSON() ([]byte, error) { ... }
```

### Key Features
- **Birthdate storage** - Date of birth stored as DATE type
- **Dynamic age calculation** - Age computed from birthdate, not stored
- **JSON compatibility** - API responses include both birthdate and computed age
- **Validation** - Business rules enforced in model layer
- **Database tags** - `db:` tags for sqlx struct mapping