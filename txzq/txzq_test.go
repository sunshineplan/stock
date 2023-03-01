package txzq

import (
	"testing"
)

func TestTXZQ(t *testing.T) {
	s := TXZQ{Index: "SSE", Code: "600309"}
	if n := s.get().Realtime.Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	s = TXZQ{Index: "SSE", Code: "688318"}
	if n := s.get().Realtime.Name; n != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", n)
	}

	s = TXZQ{Index: "SZSE", Code: "002142"}
	if n := s.get().Realtime.Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}

	s = TXZQ{Index: "SZSE", Code: "300059"}
	if n := s.get().Realtime.Name; n != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", n)
	}
}

func TestTXZQSuggests(t *testing.T) {
	s := Suggests("whhx")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	//s = Suggests("cfqs")
	//if len(s) == 0 {
	//	t.Fatal("no result")
	//}
	//if n := s[0].Name; n != "财富趋势" {
	//	t.Errorf("expected %q; got %q", "财富趋势", n)
	//}

	s = Suggests("nbyh")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}

	s = Suggests("dfcf")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", n)
	}
}
