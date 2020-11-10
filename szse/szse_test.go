package szse

import (
	"testing"
)

func TestSZSE(t *testing.T) {
	s := SZSE{Code: "002142"}
	if n := s.get().Realtime.Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}
}

func TestSuggests(t *testing.T) {
	s := Suggests("nbyh")
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}
}
