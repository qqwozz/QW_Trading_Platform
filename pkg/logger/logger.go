package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
)

type Entry struct {
	Time      string      `json:"time"`
	Level     Level       `json:"level"`
	Service   string      `json:"service"`
	Message   string      `json:"message"`
	RequestID string      `json:"request_id,omitempty"`
	Method    string      `json:"method,omitempty"`
	Path      string      `json:"path,omitempty"`
	Status    int         `json:"status,omitempty"`
	Duration  string      `json:"duration,omitempty"`
	Extra     interface{} `json:"extra,omitempty"`
}

type Logger struct {
	service string
}

func New(service string) *Logger {
	return &Logger{service: service}
}

func (l *Logger) log(level Level, msg string, fields map[string]interface{}) {
	entry := Entry{
		Time:    time.Now().UTC().Format(time.RFC3339Nano),
		Level:   level,
		Service: l.service,
		Message: msg,
	}
	if fields != nil {
		entry.Extra = fields
	}
	json.NewEncoder(os.Stdout).Encode(entry)
}

func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelInfo, msg, f)
}

func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelWarn, msg, f)
}

func (l *Logger) Error(msg string, fields ...map[string]interface{}) {
	var f map[string]interface{}
	if len(fields) > 0 {
		f = fields[0]
	}
	l.log(LevelError, msg, f)
}

func (l *Logger) Request(method, path string, status int, duration time.Duration, requestID string) {
	l.log(LevelInfo, "request", map[string]interface{}{
		"method":     method,
		"path":       path,
		"status":     status,
		"duration":   duration.String(),
		"request_id": requestID,
	})
}

func (l *Logger) Fatal(msg string, err error) {
	l.log(LevelError, msg, map[string]interface{}{
		"error": fmt.Sprintf("%v", err),
	})
	os.Exit(1)
}
