package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	logFileHandle *os.File
)

// NewLogger creates and configures a new logger instance
func NewLogger() *logrus.Logger {
	logger := logrus.New()

	// Default configuration
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)

	// Configure from environment or config if available
	configureLogger(logger)

	return logger
}

// SetLogLevel sets the logging level
func SetLogLevel(logger *logrus.Logger, level string) {
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// SetLogFormat sets the logging format
func SetLogFormat(logger *logrus.Logger, format string) {
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	default:
		// Keep current formatter
	}
}

// SetLogFile sets the log output file
func SetLogFile(logger *logrus.Logger, filename string) error {
	// Create log directory if it doesn't exist
	logDir := filepath.Dir(filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// Close previously opened file to avoid descriptor leaks
	if logFileHandle != nil {
		_ = logFileHandle.Close()
		logFileHandle = nil
	}

	// Open log file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	logFileHandle = file

	// Set output to both file and stdout
	multiWriter := io.MultiWriter(os.Stdout, file)
	logger.SetOutput(multiWriter)

	return nil
}

// CloseLogFile closes the active log file if one is open
func CloseLogFile() error {
	if logFileHandle != nil {
		err := logFileHandle.Close()
		logFileHandle = nil
		return err
	}
	return nil
}

// configureLogger applies configuration from environment variables
func configureLogger(logger *logrus.Logger) {
	// Check for log level in environment
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		SetLogLevel(logger, level)
	}

	// Check for log format in environment
	if format := os.Getenv("LOG_FORMAT"); format == "json" {
		SetLogFormat(logger, "json")
	}

	// Check for log file in environment
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		SetLogFile(logger, logFile)
	} else {
		CloseLogFile()
	}
}

// GetLogger is a placeholder for refactoring.
// It is recommended to pass logger instances instead of using a global logger.
var globalLogger *logrus.Logger

func GetLogger() *logrus.Logger {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	return globalLogger
}

// LoggerWithFields creates a logger with predefined fields
func LoggerWithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}

// ContextLogger creates a logger with context fields
type ContextLogger struct {
	*logrus.Entry
}

// NewContextLogger creates a new context logger
func NewContextLogger(context string) *ContextLogger {
	return &ContextLogger{
		Entry: GetLogger().WithField("context", context),
	}
}

// WithField adds a field to the context logger
func (cl *ContextLogger) WithField(key string, value interface{}) *ContextLogger {
	return &ContextLogger{
		Entry: cl.Entry.WithField(key, value),
	}
}

// WithFields adds multiple fields to the context logger
func (cl *ContextLogger) WithFields(fields logrus.Fields) *ContextLogger {
	return &ContextLogger{
		Entry: cl.Entry.WithFields(fields),
	}
}

// Emergency log level for critical failures
func Emergency(args ...interface{}) {
	GetLogger().WithField("severity", "emergency").Fatal(args...)
}

// Alert log level for conditions requiring immediate action
func Alert(args ...interface{}) {
	GetLogger().WithField("severity", "alert").Error(args...)
}

// Critical log level for critical conditions
func Critical(args ...interface{}) {
	GetLogger().WithField("severity", "critical").Error(args...)
}

// Notice log level for normal but significant conditions
func Notice(args ...interface{}) {
	GetLogger().WithField("severity", "notice").Info(args...)
}

// LogOperation logs the start and end of an operation
func LogOperation(operation string, fn func() error) error {
	logger := GetLogger().WithField("operation", operation)
	logger.Info("Starting operation")

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"duration": duration,
			"error":    err,
		}).Error("Operation failed")
	} else {
		logger.WithField("duration", duration).Info("Operation completed successfully")
	}

	return err
}

// LogWithDuration logs with execution duration
func LogWithDuration(message string, start time.Time, fields ...logrus.Fields) {
	duration := time.Since(start)
	entry := GetLogger().WithField("duration", duration)

	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}

	entry.Info(message)
}

// Structured logging helpers
type LogEvent struct {
	Level   logrus.Level
	Message string
	Fields  logrus.Fields
}

// LogEvents logs multiple events in a structured way
func LogEvents(events []LogEvent) {
	logger := GetLogger()
	for _, event := range events {
		entry := logger.WithFields(event.Fields)
		entry.Log(event.Level, event.Message)
	}
}

// Performance logging
type PerformanceLogger struct {
	logger *logrus.Entry
	start  time.Time
}

// NewPerformanceLogger creates a performance logger
func NewPerformanceLogger(operation string) *PerformanceLogger {
	return &PerformanceLogger{
		logger: GetLogger().WithField("operation", operation),
		start:  time.Now(),
	}
}

// Checkpoint logs a checkpoint with elapsed time
func (pl *PerformanceLogger) Checkpoint(checkpoint string) {
	elapsed := time.Since(pl.start)
	pl.logger.WithFields(logrus.Fields{
		"checkpoint": checkpoint,
		"elapsed":    elapsed,
	}).Debug("Performance checkpoint")
}

// Finish logs the completion with total duration
func (pl *PerformanceLogger) Finish(message string) {
	duration := time.Since(pl.start)
	pl.logger.WithField("total_duration", duration).Info(message)
}

// FinishWithError logs completion with error and total duration
func (pl *PerformanceLogger) FinishWithError(message string, err error) {
	duration := time.Since(pl.start)
	pl.logger.WithFields(logrus.Fields{
		"total_duration": duration,
		"error":          err,
	}).Error(message)
}
