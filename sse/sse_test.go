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

func TestSSESuggests(t *testing.T) {
	s := Suggests("whhx")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	s = Suggests("cfqs")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", n)
	}
}
