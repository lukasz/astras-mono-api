package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LocalLogger provides file-based logging for local development
type LocalLogger struct {
	*Logger
	logFile *os.File
	mutex   sync.Mutex
}

// LocalConfig holds configuration for local file logging
type LocalConfig struct {
	Config      // Embed base config
	LogDir      string // Directory for log files
	LogFileName string // Log file name (optional, defaults to service-YYYYMMDD.log)
	MaxFileSize int64  // Max file size in bytes before rotation
}

// NewLocalLogger creates a logger that writes to local files
func NewLocalLogger(config LocalConfig) (*LocalLogger, error) {
	// Set default log directory
	if config.LogDir == "" {
		pwd, _ := os.Getwd()
		config.LogDir = filepath.Join(pwd, "logs")
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Set default log file name
	if config.LogFileName == "" {
		today := time.Now().Format("2006-01-02")
		config.LogFileName = fmt.Sprintf("astras-%s-%s.log", config.ServiceName, today)
	}

	logPath := filepath.Join(config.LogDir, config.LogFileName)

	// Open log file (create or append)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create multi-writer (stdout + file)
	config.Config.Output = io.MultiWriter(os.Stdout, logFile)

	// Create base logger
	baseLogger := New(config.Config)

	return &LocalLogger{
		Logger:  baseLogger,
		logFile: logFile,
	}, nil
}

// Close closes the log file
func (ll *LocalLogger) Close() error {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()

	if ll.logFile != nil {
		return ll.logFile.Close()
	}
	return nil
}

// RotateLogFile rotates the log file (useful for long-running processes)
func (ll *LocalLogger) RotateLogFile() error {
	ll.mutex.Lock()
	defer ll.mutex.Unlock()

	if ll.logFile != nil {
		ll.logFile.Close()
	}

	// Create new log file with timestamp
	timestamp := time.Now().Format("2006-01-02-150405")
	newFileName := fmt.Sprintf("astras-%s-%s.log", ll.serviceName, timestamp)
	
	logDir := filepath.Dir(ll.logFile.Name())
	newLogPath := filepath.Join(logDir, newFileName)

	newLogFile, err := os.OpenFile(newLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to create new log file: %w", err)
	}

	ll.logFile = newLogFile

	// Update output to new file
	ll.output = io.MultiWriter(os.Stdout, newLogFile)
	ll.stdLogger.SetOutput(ll.output)

	return nil
}

// LocalDevelopmentSetup configures logging for local development
func LocalDevelopmentSetup(serviceName string) (*LocalLogger, error) {
	config := LocalConfig{
		Config: Config{
			ServiceName: serviceName,
			Environment: "local",
			MinLevel:    DEBUG, // More verbose for local development
		},
		LogDir: "logs",
	}

	logger, err := NewLocalLogger(config)
	if err != nil {
		return nil, err
	}

	// Log startup message
	logger.Info(context.Background(), fmt.Sprintf("%s started in local development mode", serviceName),
		String("log_file", logger.logFile.Name()),
		String("environment", "local"),
	)

	return logger, nil
}

// SAMLocalLogger creates a logger compatible with SAM Local
func SAMLocalLogger(serviceName string) (*Logger, error) {
	// SAM Local forwards logs to stdout, so we just use console logging
	return New(Config{
		ServiceName: serviceName,
		Environment: "sam-local",
		MinLevel:    INFO,
		Output:      os.Stdout, // SAM captures this
	}), nil
}