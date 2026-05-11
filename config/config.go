package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	URL            string
	Method         string
	Headers        map[string]string
	TotalRecords   int
	WorkerCount    int
	RequestTimeout time.Duration
	BuildPayload   func(index int) (payload any, label string)
}

type fileConfig struct {
	URL            string            `json:"url"`
	Method         string            `json:"method"`
	Headers        map[string]string `json:"headers"`
	WorkerCount    int               `json:"worker_count"`
	TotalRecords   int               `json:"total_records"`
	TimeoutSeconds int               `json:"timeout_seconds"`
}

func LoadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return Config{}, fmt.Errorf("parse config %q: %w", path, err)
	}

	applyEnvOverrides(&fc)
	applyDefaults(&fc)

	return Config{
		URL:            fc.URL,
		Method:         fc.Method,
		Headers:        fc.Headers,
		WorkerCount:    fc.WorkerCount,
		TotalRecords:   fc.TotalRecords,
		RequestTimeout: time.Duration(fc.TimeoutSeconds) * time.Second,
	}, nil
}

func overrideString(key string, dst *string) {
	if v := os.Getenv(key); v != "" {
		*dst = v
	}
}

func overrideInt(key string, dst *int) {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			*dst = n
		}
	}
}

func applyEnvOverrides(fc *fileConfig) {
	overrideString("BULKIFY_URL", &fc.URL)
	overrideString("BULKIFY_METHOD", &fc.Method)
	overrideInt("BULKIFY_WORKER_COUNT", &fc.WorkerCount)
	overrideInt("BULKIFY_TOTAL_RECORDS", &fc.TotalRecords)
	overrideInt("BULKIFY_TIMEOUT_SECONDS", &fc.TimeoutSeconds)
}

func applyDefaults(fc *fileConfig) {
	if fc.Method == "" {
		fc.Method = "POST"
	}
	if fc.WorkerCount <= 0 {
		fc.WorkerCount = 10
	}
	if fc.TotalRecords <= 0 {
		fc.TotalRecords = 1
	}
	if fc.TimeoutSeconds <= 0 {
		fc.TimeoutSeconds = 10
	}
	if fc.Headers == nil {
		fc.Headers = map[string]string{"Content-Type": "application/json"}
	}
}
