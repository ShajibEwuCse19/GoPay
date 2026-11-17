package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestJSONLogger_HasServiceAndSource(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{AddSource: true})
	l := slog.New(handler).With(slog.String("service", "FINTECH"))

	l.Info("test message", "port", 8080)

	var output map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &output); err != nil {
		t.Fatalf("log line is not valid JSON: %v", err)
	}

	// 1. service field check
	if output["service"] != "FINTECH" {
		t.Errorf("expected service 'FINTECH', got %v", output["service"])
	}

	// 2. source field check (caller file & line)
	source, ok := output["source"].(map[string]interface{})
	if !ok || source["file"] == "" {
		t.Errorf("expected caller source to be present, got %v", output["source"])
	}

	// 3. level and msg check
	if output["level"] != "INFO" || output["msg"] != "test message" {
		t.Errorf("level or msg mismatch: %v", output)
	}
}

func TestInit_Levels(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"debug"},
		{"info"},
		{"warn"},
		{"error"},
		{"invalid_fallback_to_info"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := Init("TEST_SERVICE", tt.input)
			if l == nil {
				t.Fatalf("Init returned nil logger for level %q", tt.input)
			}
		})
	}
}
