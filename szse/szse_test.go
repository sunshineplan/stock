package szse

import (
	"log"
	"testing"
)

func TestSZSE(t *testing.T) {
	s := SZSE{Code: "002142"}
	if n := s.get().Realtime.Name; n != "宁波银行" {
		log.Printf("expected %q; got %q", "宁波银行", n)
	}

	s = SZSE{Code: "300059"}
	if n := s.get().Realtime.Name; n != "东方财富" {
		log.Printf("expected %q; got %q", "东方财富", n)
	}
}

func TestSZSESuggests(t *testing.T) {
	s := Suggests("nbyh")
	if len(s) == 0 {
		log.Print("no result")
		return
	}
	if n := s[0].Name; n != "宁波银行" {
		log.Printf("expected %q; got %q", "宁波银行", n)
	}

	s = Suggests("dfcf")
	if len(s) == 0 {
		log.Print("no result")
		return
	}
	if n := s[0].Name; n != "东方财富" {
		log.Printf("expected %q; got %q", "东方财富", n)
	}
}
