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

	s = EastMoney{Index: "BSE", Code: "430047"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "诺思兰德" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", name)
	}

	s = EastMoney{Index: "BSE", Code: "834021"}
	utils.Retry(
		func() error {
			name = s.getRealtime().Realtime.Name
			if name != "流金岁月" {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if name != "流金岁月" {
		t.Errorf("expected %q; got %q", "流金岁月", name)
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
	if i := s[0].Index; i != "SSE" {
		t.Errorf("expected %q; got %q", "SSE", i)
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
	if i := s[0].Index; i != "SSE" {
		t.Errorf("expected %q; got %q", "SSE", i)
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
	if i := s[0].Index; i != "SZSE" {
		t.Errorf("expected %q; got %q", "SZSE", i)
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
	if i := s[0].Index; i != "SZSE" {
		t.Errorf("expected %q; got %q", "SZSE", i)
	}
	if n := s[0].Name; n != "东方财富" {
		t.Errorf("expected %q; got %q", "东方财富", n)
	}

	utils.Retry(
		func() error {
			s = Suggests("nsld")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "BSE" {
		t.Errorf("expected %q; got %q", "BSE", i)
	}
	if n := s[0].Name; n != "诺思兰德" {
		t.Errorf("expected %q; got %q", "诺思兰德", n)
	}

	utils.Retry(
		func() error {
			s = Suggests("ljsy")
			if len(s) == 0 {
				return errors.New("retry")
			}
			return nil
		}, 5, 20,
	)
	if len(s) == 0 {
		t.Fatal("no result")
	}
	if i := s[0].Index; i != "BSE" {
		t.Errorf("expected %q; got %q", "BSE", i)
	}
	if n := s[0].Name; n != "流金岁月" {
		t.Errorf("expected %q; got %q", "流金岁月", n)
	}
}
