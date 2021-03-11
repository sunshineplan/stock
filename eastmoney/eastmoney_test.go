package eastmoney

import (
	"testing"
)

func TestEastMoney(t *testing.T) {
	s := EastMoney{Index: "SSE", Code: "600309"}
	if n := s.getRealtime().Realtime.Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}
	s = EastMoney{Index: "SZSE", Code: "002142"}
	if n := s.getRealtime().Realtime.Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}
}

func TestSuggests(t *testing.T) {
	s := Suggests("whhx")
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}
	s = Suggests("nbyh")
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}
}
