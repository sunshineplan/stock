package sse

import (
	"testing"
)

func TestSSE(t *testing.T) {
	s := SSE{Code: "600309"}
	if n := s.getRealtime().Realtime.Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}
}

func TestSuggests(t *testing.T) {
	s := Suggests("whhx")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}
}
