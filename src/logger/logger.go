// Package logger provides a minimal, zero-dependency logging implementation
// using Go's standard library. Supports log levels, basic file output,
// and timestamp formatting.

package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// Level represents log severity
type Level int

const (
	TraceLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Logger wraps Go's standard log package with level filtering
type Logger struct {
	mu     sync.Mutex
	level  Level
	out    io.Writer
	colors bool
}

// NewLogger creates a logger with the given level and output
func NewLogger(levelStr string, logFile string, noColors bool) *Logger {
	// Parse log level
	level := InfoLevel
	switch strings.ToLower(levelStr) {
	case "trace":
		level = TraceLevel
	case "debug":
		level = DebugLevel
	case "info":
		level = InfoLevel
	case "warn":
		level = WarnLevel
	case "error":
		level = ErrorLevel
	case "fatal":
		level = FatalLevel
	}

	// Set output destination
	output := os.Stdout
	if logFile != "" && logFile != "stdout" && logFile != "/dev/stdout" {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can't open log file: %s\n", err)
			os.Exit(1)
		}
		output = f
	}

	return &Logger{
		level:  level,
		out:    output,
		colors: !noColors,
	}
}

// formatTime returns timestamp string
func formatTime() string {
	return time.Now().Format("[2006-01-02 15:04:05.000 MST]")
}

// formatFields converts key-value pairs to string
func formatFields(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}
	parts := make([]string, 0, len(fields))
	for k, v := range fields {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return " " + strings.Join(parts, " ")
}

// colorPrefix returns ANSI color codes if colors enabled
func (l *Logger) colorPrefix(level string) string {
	if !l.colors {
		return "[" + level + "]"
	}

	colors := map[string]string{
		"TRACE": "\033[36mTRACE\033[0m", // cyan
		"DEBUG": "\033[36mDEBUG\033[0m", // cyan
		"INFO":  "\033[32mINFO \033[0m", // green
		"WARN":  "\033[33mWARN \033[0m", // yellow
		"ERROR": "\033[31mERROR\033[0m", // red
		"FATAL": "\033[31mFATAL\033[0m", // red
	}

	if c, ok := colors[level]; ok {
		return c
	}
	return "[" + level + "]"
}

// Trace logs a message at trace level
func (l *Logger) Trace(msg string, fields map[string]interface{}) {
	if l.level <= TraceLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s %s%s\n", l.colorPrefix("TRACE"), formatTime(), msg, formatFields(fields))
	}
}

// Debug logs a message at debug level
func (l *Logger) Debug(msg string, fields map[string]interface{}) {
	if l.level <= DebugLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s %s%s\n", l.colorPrefix("DEBUG"), formatTime(), msg, formatFields(fields))
	}
}

// Info logs a message at info level
func (l *Logger) Info(msg string, fields map[string]interface{}) {
	if l.level <= InfoLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s %s%s\n", l.colorPrefix("INFO"), formatTime(), msg, formatFields(fields))
	}
}

// Warn logs a message at warn level
func (l *Logger) Warn(msg string, fields map[string]interface{}) {
	if l.level <= WarnLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s %s%s\n", l.colorPrefix("WARN"), formatTime(), msg, formatFields(fields))
	}
}

// Error logs a message at error level
func (l *Logger) Error(msg string, fields map[string]interface{}) {
	if l.level <= ErrorLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s %s%s\n", l.colorPrefix("ERROR"), formatTime(), msg, formatFields(fields))
	}
}

// Fatal logs a message at fatal level and exits
func (l *Logger) Fatal(msg string, fields map[string]interface{}) {
	if l.level <= FatalLevel {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Fprintf(l.out, "%s %s FATAL: %s%s\n", l.colorPrefix("FATAL"), formatTime(), msg, formatFields(fields))
		os.Exit(1)
	}
}
