//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

const (
	GOARCH = "amd64"
	GOOS   = "linux"
)

var services = []string{
	"kid-service",
	"caregiver-service",
	"star-service",
	"migration-service",
}

// Build all services for deployment
func (Build) All() error {
	fmt.Println("Building all services...")

	for _, service := range services {
		if err := buildService(service); err != nil {
			return fmt.Errorf("failed to build %s: %w", service, err)
		}
	}

	fmt.Println("All services built successfully!")
	return nil
}

// Build a specific service
func (Build) Service(service string) error {
	fmt.Printf("Building service: %s\n", service)
	return buildService(service)
}

// Build Kid service
func (Build) Kid() error {
	return buildService("kid-service")
}

// Build Caregiver service
func (Build) Caregiver() error {
	return buildService("caregiver-service")
}

// Build Star service
func (Build) Star() error {
	return buildService("star-service")
}

// Build Migration service
func (Build) Migration() error {
	return buildService("migration-service")
}

// Build Kid service for local development
func (Build) KidLocal() error {
	return buildServiceLocal("kid-service")
}

// Build Caregiver service for local development
func (Build) CaregiverLocal() error {
	return buildServiceLocal("caregiver-service")
}

// Build Star service for local development
func (Build) StarLocal() error {
	return buildServiceLocal("star-service")
}

// Build Migration service for local development
func (Build) MigrationLocal() error {
	return buildServiceLocal("migration-service")
}

// Build all services for local development
func (Build) AllLocal() error {
	fmt.Println("Building all services for local development...")

	for _, service := range services {
		if err := buildServiceLocal(service); err != nil {
			return fmt.Errorf("failed to build %s locally: %w", service, err)
		}
	}

	fmt.Println("All services built successfully for local development!")
	return nil
}

func buildService(service string) error {
	servicePath := filepath.Join("cmd", service)
	outputPath := filepath.Join("bin", service, "bootstrap")

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	env := map[string]string{
		"GOOS":        GOOS,
		"GOARCH":      GOARCH,
		"CGO_ENABLED": "0",
	}

	fmt.Printf("Building %s -> %s\n", servicePath, outputPath)
	return sh.RunWithV(env, "go", "build", "-ldflags", "-s -w", "-o", outputPath, "./"+servicePath)
}

func buildServiceLocal(service string) error {
	servicePath := filepath.Join("cmd", service)
	outputPath := filepath.Join("bin", service, "bootstrap")

	// Create output directory
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build for local OS (no cross-compilation)
	env := map[string]string{
		"CGO_ENABLED": "0",
	}

	fmt.Printf("Building %s for local development -> %s\n", servicePath, outputPath)
	return sh.RunWithV(env, "go", "build", "-ldflags", "-s -w", "-o", outputPath, "./"+servicePath)
}

type Deploy mg.Namespace

// Deploy all services
func (Deploy) All() error {
	mg.Deps(Build.All)
	fmt.Println("Deploying all services...")

	for _, service := range services {
		if err := deployService(service); err != nil {
			return fmt.Errorf("failed to deploy %s: %w", service, err)
		}
	}

	fmt.Println("All services deployed successfully!")
	return nil
}

// Deploy a specific service
func (Deploy) Service(service string) error {
	mg.Deps(mg.F(buildService, service))
	return deployService(service)
}

// Deploy Kid service
func (Deploy) Kid() error {
	mg.Deps(Build.Kid)
	return deployService("kid-service")
}

// Deploy Caregiver service
func (Deploy) Caregiver() error {
	mg.Deps(Build.Caregiver)
	return deployService("caregiver-service")
}

// Deploy Star service
func (Deploy) Star() error {
	mg.Deps(Build.Star)
	return deployService("star-service")
}

// Deploy Migration service
func (Deploy) Migration() error {
	mg.Deps(Build.Migration)
	return deployService("migration-service")
}

// Deploy infrastructure only
func (Deploy) Infrastructure() error {
	fmt.Println("Deploying infrastructure...")
	
	cmd := exec.Command("serverless", "deploy", "--verbose")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "SLS_CONFIG_FILE=serverless-infrastructure.yml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Deploy logging infrastructure
func (Deploy) Logging() error {
	fmt.Println("Deploying logging infrastructure...")
	
	cmd := exec.Command("serverless", "deploy", "--verbose")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "SLS_CONFIG_FILE=serverless-logging.yml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func deployService(service string) error {
	servicePath := filepath.Join("services", service)
	fmt.Printf("Deploying service: %s from %s\n", service, servicePath)

	cmd := exec.Command("serverless", "deploy", "--verbose")
	cmd.Dir = servicePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

type Migration mg.Namespace

// Deploy infrastructure and migration service
func (Migration) Deploy() error {
	fmt.Println("Deploying database infrastructure and migration service...")
	
	// First deploy infrastructure
	if err := (Deploy{}).Infrastructure(); err != nil {
		return fmt.Errorf("failed to deploy infrastructure: %w", err)
	}
	
	// Then deploy migration service
	if err := (Deploy{}).Migration(); err != nil {
		return fmt.Errorf("failed to deploy migration service: %w", err)
	}
	
	return nil
}

// Run migrations via deployed Lambda function
func (Migration) Migrate() error {
	fmt.Println("Running database migrations...")
	// This would call the deployed Lambda function
	// For now, we'll use local execution
	return sh.RunV("curl", "-X", "POST", "-H", "Content-Type: application/json", 
		"https://api.gateway.url/migrations/migrate")
}

// Rollback migrations via deployed Lambda function
func (Migration) Rollback() error {
	fmt.Println("Rolling back database migrations...")
	// This would call the deployed Lambda function
	return sh.RunV("curl", "-X", "POST", "-H", "Content-Type: application/json",
		"https://api.gateway.url/migrations/rollback")
}

// Check migration status
func (Migration) Status() error {
	fmt.Println("Checking migration status...")
	// This would call the deployed Lambda function
	return sh.RunV("curl", "-X", "GET", "-H", "Content-Type: application/json",
		"https://api.gateway.url/migrations/status")
}

type Test mg.Namespace

// Run tests for all services
func (Test) All() error {
	fmt.Println("Running tests for all services...")
	return sh.RunV("go", "test", "./...")
}

// Run tests with coverage
func (Test) Coverage() error {
	fmt.Println("Running tests with coverage...")
	return sh.RunV("go", "test", "-cover", "./...")
}

type Clean mg.Namespace

// Clean build artifacts
func (Clean) Build() error {
	fmt.Println("Cleaning build artifacts...")
	return os.RemoveAll("bin")
}

// Clean Serverless artifacts
func (Clean) Deploy() error {
	fmt.Println("Cleaning Serverless artifacts...")

	for _, service := range services {
		servicePath := filepath.Join("services", service, ".serverless")
		if err := os.RemoveAll(servicePath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to clean %s: %w", servicePath, err)
		}
	}

	return nil
}

// Clean everything
func (Clean) All() error {
	mg.Deps(Clean.Build, Clean.Deploy)
	return nil
}

// Format Go code
func Format() error {
	fmt.Println("Formatting Go code...")
	return sh.RunV("go", "fmt", "./...")
}

// Lint Go code
func Lint() error {
	fmt.Println("Linting Go code...")
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}

	// Run golangci-lint if available
	if _, err := exec.LookPath("golangci-lint"); err == nil {
		return sh.RunV("golangci-lint", "run")
	}

	fmt.Println("golangci-lint not found, skipping advanced linting")
	return nil
}

// Tidy Go modules
func Tidy() error {
	fmt.Println("Tidying Go modules...")
	return sh.RunV("go", "mod", "tidy")
}

// Install Mage
func InstallMage() error {
	fmt.Println("Installing Mage...")
	return sh.RunV("go", "install", "github.com/magefile/mage@latest")
}

// List available services
func Services() error {
	fmt.Println("Available services:")
	for _, service := range services {
		fmt.Printf("  - %s\n", service)
	}
	return nil
}

// Default target
func Default() {
	fmt.Println("Available targets:")
	fmt.Println("  mage build:all        - Build all services (for AWS Lambda)")
	fmt.Println("  mage build:allLocal   - Build all services (for local development)")
	fmt.Println("  mage build:kid        - Build kid service")
	fmt.Println("  mage build:kidLocal   - Build kid service (for local development)")
	fmt.Println("  mage build:caregiver  - Build caregiver service")
	fmt.Println("  mage build:caregiverLocal - Build caregiver service (for local development)")
	fmt.Println("  mage build:star       - Build star service")
	fmt.Println("  mage build:starLocal  - Build star service (for local development)")
	fmt.Println("  mage deploy:all       - Deploy all services")
	fmt.Println("  mage deploy:kid       - Deploy kid service")
	fmt.Println("  mage deploy:caregiver - Deploy caregiver service")
	fmt.Println("  mage deploy:star      - Deploy star service")
	fmt.Println("  mage test:all         - Run all tests")
	fmt.Println("  mage test:coverage    - Run tests with coverage")
	fmt.Println("  mage clean:all        - Clean all artifacts")
	fmt.Println("  mage format           - Format Go code")
	fmt.Println("  mage lint             - Lint Go code")
	fmt.Println("  mage tidy             - Tidy Go modules")
	fmt.Println("  mage services         - List available services")
}
