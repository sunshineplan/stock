package bse

import (
	"testing"
)

func TestBSE(t *testing.T) {
	s := BSE{Code: "920047"}
	if n := s.getRealtime().Realtime.Name; n != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", n)
	}

	s = BSE{Code: "920185"}
	if n := s.getRealtime().Realtime.Name; n != "贝特瑞" {
		t.Errorf("expected %q; got %q", "贝特瑞", n)
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

	s = Suggests("btr")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "贝特瑞" {
		t.Errorf("expected %q; got %q", "贝特瑞", n)
	}
}
