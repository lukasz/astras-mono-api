package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// LogLevel represents the severity of the log message
type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp    time.Time              `json:"@timestamp"`
	Level        LogLevel               `json:"level"`
	Service      string                 `json:"service"`
	Message      string                 `json:"message"`
	RequestID    string                 `json:"request_id,omitempty"`
	UserID       string                 `json:"user_id,omitempty"`
	Operation    string                 `json:"operation,omitempty"`
	Duration     *int64                 `json:"duration,omitempty"` // milliseconds
	StatusCode   *int                   `json:"status_code,omitempty"`
	Error        string                 `json:"error,omitempty"`
	File         string                 `json:"file,omitempty"`
	Line         int                    `json:"line,omitempty"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
	Environment  string                 `json:"environment"`
	AWSRequestID string                 `json:"aws_request_id,omitempty"`
	Version      string                 `json:"version,omitempty"`
}

// Logger provides structured logging functionality
type Logger struct {
	serviceName string
	environment string
	version     string
	minLevel    LogLevel
	output      io.Writer
	stdLogger   *log.Logger
}

// Config holds logger configuration
type Config struct {
	ServiceName string
	Environment string
	Version     string
	MinLevel    LogLevel
	Output      io.Writer
}

// New creates a new structured logger
func New(config Config) *Logger {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	
	if config.MinLevel == "" {
		config.MinLevel = INFO
	}

	if config.Environment == "" {
		config.Environment = os.Getenv("STAGE")
		if config.Environment == "" {
			config.Environment = "dev"
		}
	}

	return &Logger{
		serviceName: config.ServiceName,
		environment: config.Environment,
		version:     config.Version,
		minLevel:    config.MinLevel,
		output:      config.Output,
		stdLogger:   log.New(config.Output, "", 0),
	}
}

// shouldLog determines if a message should be logged based on level
func (l *Logger) shouldLog(level LogLevel) bool {
	levelOrder := map[LogLevel]int{
		DEBUG: 0,
		INFO:  1,
		WARN:  2,
		ERROR: 3,
	}
	return levelOrder[level] >= levelOrder[l.minLevel]
}

// log writes a structured log entry
func (l *Logger) log(ctx context.Context, level LogLevel, message string, fields ...Field) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp:   time.Now().UTC(),
		Level:       level,
		Service:     l.serviceName,
		Message:     message,
		Environment: l.environment,
		Version:     l.version,
		Extra:       make(map[string]interface{}),
	}

	// Add AWS Lambda context if available
	if lambdaCtx, ok := lambdacontext.FromContext(ctx); ok {
		entry.AWSRequestID = lambdaCtx.AwsRequestID
	}

	// Add file and line information for ERROR level
	if level == ERROR {
		if pc, file, line, ok := runtime.Caller(2); ok {
			entry.File = trimPath(file)
			entry.Line = line
			
			// Try to get function name
			if fn := runtime.FuncForPC(pc); fn != nil {
				funcName := fn.Name()
				if idx := strings.LastIndex(funcName, "/"); idx >= 0 {
					funcName = funcName[idx+1:]
				}
				entry.Extra["function"] = funcName
			}
		}
	}

	// Apply fields
	for _, field := range fields {
		field.apply(&entry)
	}

	// Convert to JSON and write
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to standard logging if JSON marshaling fails
		l.stdLogger.Printf("[%s] %s: %s (JSON marshal error: %v)", level, l.serviceName, message, err)
		return
	}

	l.stdLogger.Println(string(jsonData))
}

// Debug logs a debug message
func (l *Logger) Debug(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(ctx context.Context, message string, fields ...Field) {
	l.log(ctx, ERROR, message, fields...)
}

// Field represents a structured logging field
type Field interface {
	apply(*LogEntry)
}

// field implementations
type stringField struct {
	key   string
	value string
}

func (f stringField) apply(entry *LogEntry) {
	switch f.key {
	case "request_id":
		entry.RequestID = f.value
	case "user_id":
		entry.UserID = f.value
	case "operation":
		entry.Operation = f.value
	case "error":
		entry.Error = f.value
	default:
		if entry.Extra == nil {
			entry.Extra = make(map[string]interface{})
		}
		entry.Extra[f.key] = f.value
	}
}

type intField struct {
	key   string
	value int
}

func (f intField) apply(entry *LogEntry) {
	switch f.key {
	case "status_code":
		entry.StatusCode = &f.value
	default:
		if entry.Extra == nil {
			entry.Extra = make(map[string]interface{})
		}
		entry.Extra[f.key] = f.value
	}
}

type int64Field struct {
	key   string
	value int64
}

func (f int64Field) apply(entry *LogEntry) {
	switch f.key {
	case "duration":
		entry.Duration = &f.value
	default:
		if entry.Extra == nil {
			entry.Extra = make(map[string]interface{})
		}
		entry.Extra[f.key] = f.value
	}
}

type interfaceField struct {
	key   string
	value interface{}
}

func (f interfaceField) apply(entry *LogEntry) {
	if entry.Extra == nil {
		entry.Extra = make(map[string]interface{})
	}
	entry.Extra[f.key] = f.value
}

// Field constructors
func String(key, value string) Field {
	return stringField{key, value}
}

func Int(key string, value int) Field {
	return intField{key, value}
}

func Int64(key string, value int64) Field {
	return int64Field{key, value}
}

func Duration(duration time.Duration) Field {
	return int64Field{"duration", duration.Milliseconds()}
}

func Error(err error) Field {
	if err == nil {
		return stringField{"error", ""}
	}
	return stringField{"error", err.Error()}
}

func Any(key string, value interface{}) Field {
	return interfaceField{key, value}
}

// Helper fields
func RequestID(id string) Field {
	return String("request_id", id)
}

func UserID(id string) Field {
	return String("user_id", id)
}

func Operation(op string) Field {
	return String("operation", op)
}

func StatusCode(code int) Field {
	return Int("status_code", code)
}

// Database operation logging helpers
func DatabaseOperation(operation, table string, duration time.Duration, err error) []Field {
	fields := []Field{
		String("db_operation", operation),
		String("db_table", table),
		Duration(duration),
	}
	
	if err != nil {
		fields = append(fields, Error(err))
	}
	
	return fields
}

// HTTP request logging helpers
func HTTPRequest(method, path string, statusCode int, duration time.Duration) []Field {
	return []Field{
		String("http_method", method),
		String("http_path", path),
		Int("status_code", statusCode),
		Duration(duration),
	}
}

// trimPath removes the GOPATH/GOROOT prefix from file paths for cleaner logs
func trimPath(path string) string {
	// Try to find common prefixes to remove
	if idx := strings.LastIndex(path, "/internal/"); idx >= 0 {
		return path[idx+1:]
	}
	if idx := strings.LastIndex(path, "/cmd/"); idx >= 0 {
		return path[idx+1:]
	}
	if idx := strings.LastIndex(path, "/pkg/"); idx >= 0 {
		return path[idx+1:]
	}
	
	// If no common prefix found, just return the filename
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}
	
	return path
}

// DatabaseLogger is a specialized logger for database operations
type DatabaseLogger struct {
	logger *Logger
}

// NewDatabaseLogger creates a logger specialized for database operations
func NewDatabaseLogger(serviceName string) *DatabaseLogger {
	return &DatabaseLogger{
		logger: New(Config{
			ServiceName: serviceName + "-db",
			MinLevel:    INFO,
		}),
	}
}

// LogQuery logs a database query with timing and result information
func (dl *DatabaseLogger) LogQuery(ctx context.Context, query string, args []interface{}, duration time.Duration, err error) {
	level := INFO
	if err != nil {
		level = ERROR
	}

	fields := []Field{
		String("operation", "query"),
		String("query", query),
		Duration(duration),
		Any("query_args", args),
	}

	if err != nil {
		fields = append(fields, Error(err))
	}

	message := fmt.Sprintf("Database query executed in %dms", duration.Milliseconds())
	if err != nil {
		message = fmt.Sprintf("Database query failed after %dms", duration.Milliseconds())
	}

	dl.logger.log(ctx, level, message, fields...)
}