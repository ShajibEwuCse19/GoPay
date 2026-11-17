package utils

import (
	"context"
	"strings"
	"testing"
)

func TestContextRequestID(t *testing.T) {
	t.Run("nil context returns empty string", func(t *testing.T) {
		if got := GetRequestID(nil); got != "" {
			t.Errorf("GetRequestID(nil) = %q, want empty", got)
		}
	})

	t.Run("empty context returns empty string", func(t *testing.T) {
		ctx := context.Background()
		if got := GetRequestID(ctx); got != "" {
			t.Errorf("GetRequestID(empty) = %q, want empty", got)
		}
	})

	t.Run("stores and retrieves request ID", func(t *testing.T) {
		ctx := context.Background()
		reqID := "test-request-id-12345"
		ctx = WithRequestID(ctx, reqID)

		if got := GetRequestID(ctx); got != reqID {
			t.Errorf("GetRequestID = %q, want %q", got, reqID)
		}
	})

	t.Run("generates valid UUID v4 format", func(t *testing.T) {
		id1 := GenerateRequestID()
		id2 := GenerateRequestID()

		if id1 == id2 {
			t.Errorf("expected generated IDs to be unique, got duplicate %q", id1)
		}

		parts := strings.Split(id1, "-")
		if len(parts) != 5 {
			t.Errorf("invalid UUID format: %q", id1)
		}
		if len(parts[0]) != 8 || len(parts[1]) != 4 || len(parts[2]) != 4 || len(parts[3]) != 4 || len(parts[4]) != 12 {
			t.Errorf("invalid UUID component lengths in: %q", id1)
		}
	})
}
