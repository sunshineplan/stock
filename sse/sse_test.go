package sse

import (
	"testing"
)

func TestSSE(t *testing.T) {
	s := SSE{Code: "600309"}
	if s.getRealtime().Realtime.Name != "万华化学" {
		t.Error("Get sse stock error")
	}
}
