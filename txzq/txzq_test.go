package txzq

import (
	"testing"
)

func TestTXZQ(t *testing.T) {
	s := TXZQ{Index: "SSE", Code: "600309"}
	if s.get().Realtime.Name != "万华化学" {
		t.Error("Get sse stock error")
	}
	s = TXZQ{Index: "SZSE", Code: "002142"}
	if s.get().Realtime.Name != "宁波银行" {
		t.Error("Get szse stock error")
	}
}
