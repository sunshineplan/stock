package sse

import (
	"testing"
)

func TestSSE(t *testing.T) {
	s := SSE{Code: "600309"}
	if n := s.getRealtime().Realtime.Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	s = SSE{Code: "688318"}
	if n := s.getRealtime().Realtime.Name; n != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", n)
	}
}
