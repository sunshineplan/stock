package txzq

import (
	"testing"
)

func TestTXZQ(t *testing.T) {
	s := TXZQ{Index: "SSE", Code: "600309"}
	if n := s.get().Realtime.Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}
	s = TXZQ{Index: "SZSE", Code: "002142"}
	if n := s.get().Realtime.Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
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

	s = Suggests("nbyh")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}
}
