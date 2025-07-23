package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// SimpleErrorLogger provides a simple implementation of ErrorLogger
type SimpleErrorLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewSimpleErrorLogger creates a new simple error logger
func NewSimpleErrorLogger(logFile string) (*SimpleErrorLogger, error) {
	var file *os.File
	var err error

	if logFile != "" {
		file, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
	}

	logger := log.New(os.Stdout, "", 0)
	if file != nil {
		logger = log.New(file, "", 0)
	}

	return &SimpleErrorLogger{
		logger: logger,
		file:   file,
	}, nil
}

// LogError logs an error with the specified level and context
func (sel *SimpleErrorLogger) LogError(level string, message string, err error, context interface{}) {
	logEntry := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"level":     level,
		"message":   message,
		"context":   context,
	}

	if err != nil {
		logEntry["error"] = err.Error()
	}

	// Convert to JSON for structured logging
	jsonData, jsonErr := json.Marshal(logEntry)
	if jsonErr != nil {
		// Fallback to simple text logging
		sel.logger.Printf("[%s] %s: %v (context: %v)", level, message, err, context)
		return
	}

	sel.logger.Println(string(jsonData))
}

// Close closes the log file if it was opened
func (sel *SimpleErrorLogger) Close() error {
	if sel.file != nil {
		return sel.file.Close()
	}
	return nil
}

// ConsoleErrorLogger logs errors to the console only
type ConsoleErrorLogger struct{}

// NewConsoleErrorLogger creates a new console error logger
func NewConsoleErrorLogger() *ConsoleErrorLogger {
	return &ConsoleErrorLogger{}
}

// LogError logs an error to the console
func (cel *ConsoleErrorLogger) LogError(level string, message string, err error, context interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		fmt.Printf("[%s] %s - %s: %v\n", timestamp, level, message, err)
	} else {
		fmt.Printf("[%s] %s - %s\n", timestamp, level, message)
	}

	if context != nil {
		if contextJSON, jsonErr := json.MarshalIndent(context, "  ", "  "); jsonErr == nil {
			fmt.Printf("  Context: %s\n", string(contextJSON))
		} else {
			fmt.Printf("  Context: %v\n", context)
		}
	}
}
