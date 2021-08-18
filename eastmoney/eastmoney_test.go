package eastmoney

import (
	"errors"
	"testing"

	"github.com/sunshineplan/stock"
	"github.com/sunshineplan/utils"
)

func TestEastMoney(t *testing.T) {
	var name string
	s := EastMoney{Index: "SSE", Code: "600309"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "万华化学" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", name)
	}

	s = EastMoney{Index: "SSE", Code: "688318"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "财富趋势" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", name)
	}

	s = EastMoney{Index: "SZSE", Code: "002142"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "宁波银行" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", name)
	}

	s = EastMoney{Index: "SZSE", Code: "300059"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "东方财富" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", name)
	}
}

func TestSuggests(t *testing.T) {
	var s []stock.Suggest
	utils.Retry(
		func() error {
			s = Suggests("whhx")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "万华化学" {
		t.Errorf("expected %q; got %q", "万华化学", n)
	}

	utils.Retry(
		func() error {
			s = Suggests("cfqs")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "财富趋势" {
		t.Errorf("expected %q; got %q", "财富趋势", n)
	}

	utils.Retry(
		func() error {
			s = Suggests("nbyh")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "宁波银行" {
		t.Errorf("expected %q; got %q", "宁波银行", n)
	}

	utils.Retry(
		func() error {
			s = Suggests("dfcf")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if n := s[0].Name; n != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", n)
	}
}
