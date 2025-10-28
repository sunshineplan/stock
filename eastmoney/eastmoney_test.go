package eastmoney

import (
	"testing"
)

func testEastMoney(t *testing.T) {
	var name string
	s := EastMoney{Index: "SSE", Code: "600309"}
	if name = s.getRealtime().Realtime.Name; name != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", name)
	}

	s = EastMoney{Index: "SSE", Code: "688318"}
	if name = s.getRealtime().Realtime.Name; name != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", name)
	}

	s = EastMoney{Index: "SZSE", Code: "002142"}
	if name = s.getRealtime().Realtime.Name; name != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", name)
	}

	s = EastMoney{Index: "SZSE", Code: "300059"}
	if name = s.getRealtime().Realtime.Name; name != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", name)
	}

	s = EastMoney{Index: "BSE", Code: "920047"}
	if name = s.getRealtime().Realtime.Name; name != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", name)
	}

	s = EastMoney{Index: "BSE", Code: "920185"}
	if name = s.getRealtime().Realtime.Name; name != "贝特瑞" {
		t.Errorf("expected %q; got %q", "贝特瑞", name)
	}
}

func TestEastMoneySuggests(t *testing.T) {
	s := Suggests("whhx")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "SSE" {
		t.Errorf("expected %q; got %q", "SSE", i)
	}
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	s = Suggests("cfqs")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "SSE" {
		t.Errorf("expected %q; got %q", "SSE", i)
	}
	if n := s[0].Name; n != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", n)
	}

	s = Suggests("nbyh")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "SZSE" {
		t.Errorf("expected %q; got %q", "SZSE", i)
	}
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}

	s = Suggests("dfcf")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "SZSE" {
		t.Errorf("expected %q; got %q", "SZSE", i)
	}
	if n := s[0].Name; n != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", n)
	}

	s = Suggests("nsld")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "BSE" {
		t.Errorf("expected %q; got %q", "BSE", i)
	}
	if n := s[0].Name; n != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", n)
	}

	s = Suggests("btr")
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "BSE" {
		t.Errorf("expected %q; got %q", "BSE", i)
	}
	if n := s[0].Name; n != "贝特瑞" {
		t.Errorf("expected %q; got %q", "贝特瑞", n)
	}
}
