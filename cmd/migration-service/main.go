package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

const (
	MigrationsTable = "schema_migrations"
)

type MigrationService struct {
	db *sql.DB
}

type MigrationFile struct {
	Version   string
	Name      string
	Direction string // "up" or "down"
	Content   string
	FilePath  string
}

type MigrationRequest struct {
	Command string `json:"command"` // "migrate", "rollback", "status"
	Steps   int    `json:"steps,omitempty"`
}

type MigrationResponse struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	AppliedVersions []string `json:"applied_versions,omitempty"`
	Error           string   `json:"error,omitempty"`
}

func NewMigrationService() (*MigrationService, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	if dbSSLMode == "" {
		dbSSLMode = "require"
	}

	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbHost, dbPort, dbName, dbUser, dbPassword, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	ms := &MigrationService{db: db}

	// Ensure migrations table exists
	if err := ms.createMigrationsTable(); err != nil {
		return nil, fmt.Errorf("failed to create migrations table: %w", err)
	}

	return ms, nil
}

func (ms *MigrationService) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + MigrationsTable + ` (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := ms.db.Exec(query)
	return err
}

func (ms *MigrationService) loadMigrationFiles(migrationsPath string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	err := filepath.Walk(migrationsPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}

		filename := info.Name()
		parts := strings.Split(filename, ".")

		if len(parts) < 3 {
			return nil // Skip files that don't match pattern: version.name.direction.sql
		}

		version := parts[0]
		direction := parts[len(parts)-2] // second to last part

		if direction != "up" && direction != "down" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		migration := MigrationFile{
			Version:   version,
			Name:      strings.Join(parts[1:len(parts)-2], "."),
			Direction: direction,
			Content:   string(content),
			FilePath:  path,
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func (ms *MigrationService) getAppliedVersions() (map[string]time.Time, error) {
	query := `SELECT version, applied_at FROM ` + MigrationsTable + ` ORDER BY version`
	rows, err := ms.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]time.Time)
	for rows.Next() {
		var version string
		var appliedAt time.Time
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = appliedAt
	}

	return applied, rows.Err()
}

func (ms *MigrationService) migrate(migrationsPath string) error {
	migrations, err := ms.loadMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	applied, err := ms.getAppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	upMigrations := make(map[string]MigrationFile)
	for _, m := range migrations {
		if m.Direction == "up" {
			upMigrations[m.Version] = m
		}
	}

	// Sort versions for consistent ordering
	var versions []string
	for version := range upMigrations {
		if _, isApplied := applied[version]; !isApplied {
			versions = append(versions, version)
		}
	}
	sort.Strings(versions)

	if len(versions) == 0 {
		log.Println("No new migrations to apply")
		return nil
	}

	for _, version := range versions {
		migration := upMigrations[version]
		log.Printf("Applying migration %s: %s", migration.Version, migration.Name)

		tx, err := ms.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %s: %w", version, err)
		}

		// Execute migration
		if _, err := tx.Exec(migration.Content); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", version, err)
		}

		// Record migration as applied
		if _, err := tx.Exec(`INSERT INTO `+MigrationsTable+` (version) VALUES ($1)`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", version, err)
		}

		log.Printf("Successfully applied migration %s", version)
	}

	return nil
}

func (ms *MigrationService) rollback(migrationsPath string, steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps must be greater than 0")
	}

	migrations, err := ms.loadMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migration files: %w", err)
	}

	applied, err := ms.getAppliedVersions()
	if err != nil {
		return fmt.Errorf("failed to get applied versions: %w", err)
	}

	// Create map of down migrations
	downMigrations := make(map[string]MigrationFile)
	for _, m := range migrations {
		if m.Direction == "down" {
			downMigrations[m.Version] = m
		}
	}

	// Get applied versions in reverse order
	var appliedVersions []string
	for version := range applied {
		appliedVersions = append(appliedVersions, version)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(appliedVersions)))

	if len(appliedVersions) == 0 {
		log.Println("No migrations to rollback")
		return nil
	}

	// Limit to requested steps
	if steps > len(appliedVersions) {
		steps = len(appliedVersions)
	}

	for i := 0; i < steps; i++ {
		version := appliedVersions[i]
		migration, exists := downMigrations[version]
		if !exists {
			return fmt.Errorf("down migration not found for version %s", version)
		}

		log.Printf("Rolling back migration %s: %s", migration.Version, migration.Name)

		tx, err := ms.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for rollback %s: %w", version, err)
		}

		// Execute rollback
		if _, err := tx.Exec(migration.Content); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute rollback %s: %w", version, err)
		}

		// Remove migration record
		if _, err := tx.Exec(`DELETE FROM `+MigrationsTable+` WHERE version = $1`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to remove migration record %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit rollback %s: %w", version, err)
		}

		log.Printf("Successfully rolled back migration %s", version)
	}

	return nil
}

func (ms *MigrationService) status() ([]string, error) {
	applied, err := ms.getAppliedVersions()
	if err != nil {
		return nil, err
	}

	var versions []string
	for version := range applied {
		versions = append(versions, version)
	}
	sort.Strings(versions)

	return versions, nil
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received migration request: %s", request.Body)

	ms, err := NewMigrationService()
	if err != nil {
		log.Printf("Failed to create migration service: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"success": false, "error": "Failed to initialize migration service: %v"}`, err),
		}, nil
	}
	defer ms.db.Close()

	migrationsPath := "/opt/migrations" // Lambda layer path

	var response MigrationResponse

	// Parse command from path or body
	command := strings.TrimPrefix(request.Path, "/migrations/")
	if command == "" || command == request.Path {
		command = "status" // default command
	}

	switch command {
	case "migrate":
		if err := ms.migrate(migrationsPath); err != nil {
			response = MigrationResponse{
				Success: false,
				Message: "Migration failed",
				Error:   err.Error(),
			}
		} else {
			versions, _ := ms.status()
			response = MigrationResponse{
				Success:         true,
				Message:         "Migrations applied successfully",
				AppliedVersions: versions,
			}
		}

	case "rollback":
		steps := 1 // default to 1 step
		if err := ms.rollback(migrationsPath, steps); err != nil {
			response = MigrationResponse{
				Success: false,
				Message: "Rollback failed",
				Error:   err.Error(),
			}
		} else {
			versions, _ := ms.status()
			response = MigrationResponse{
				Success:         true,
				Message:         "Rollback completed successfully",
				AppliedVersions: versions,
			}
		}

	case "status":
		versions, err := ms.status()
		if err != nil {
			response = MigrationResponse{
				Success: false,
				Message: "Failed to get migration status",
				Error:   err.Error(),
			}
		} else {
			response = MigrationResponse{
				Success:         true,
				Message:         "Migration status retrieved successfully",
				AppliedVersions: versions,
			}
		}

	default:
		response = MigrationResponse{
			Success: false,
			Message: "Unknown command",
			Error:   fmt.Sprintf("Unknown command: %s", command),
		}
	}

	responseBody := fmt.Sprintf(`{"success": %t, "message": "%s"`, response.Success, response.Message)
	
	if len(response.AppliedVersions) > 0 {
		responseBody += fmt.Sprintf(`, "applied_versions": ["%s"]`, strings.Join(response.AppliedVersions, `", "`))
	}
	
	if response.Error != "" {
		responseBody += fmt.Sprintf(`, "error": "%s"`, response.Error)
	}
	
	responseBody += "}"

	statusCode := 200
	if !response.Success {
		statusCode = 400
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: responseBody,
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}