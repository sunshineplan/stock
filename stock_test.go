package main

import "testing"

func TestSSE(t *testing.T) {
	s := sse{code: "600309"}
	s.getRealtime()
	if s.name != "万华化学" {
		t.Error("Get sse stock error")
	}
}

func TestSZSE(t *testing.T) {
	s := szse{code: "002142"}
	s.getRealtime()
	if s.name != "宁波银行" {
		t.Error("Get szse stock error")
	}
}
