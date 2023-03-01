package bse

import (
	"testing"
)

func TestBSE(t *testing.T) {
	s := BSE{Code: "430047"}
	if n := s.getRealtime().Realtime.Name; n != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", n)
	}

	s = BSE{Code: "834021"}
	if n := s.getRealtime().Realtime.Name; n != "流金岁月" {
		t.Errorf("expected %q; got %q", "流金岁月", n)
	}
}

func TestBSESuggests(t *testing.T) {
	s := Suggests("nsld")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", n)
	}

	s = Suggests("ljsy")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "流金岁月" {
		t.Errorf("expected %q; got %q", "流金岁月", n)
	}
}
