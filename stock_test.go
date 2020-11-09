package main

import (
	"testing"
)

func TestSSE(t *testing.T) {
	s := sse{Code: "600309"}
	s.getRealtime()
	if s.Realtime.Name != "万华化学" {
		t.Error("Get sse stock error")
	}
}

func TestSZSE(t *testing.T) {
	s := szse{Code: "002142"}
	s.getRealtime()
	if s.Realtime.Name != "宁波银行" {
		t.Error("Get szse stock error")
	}
}
