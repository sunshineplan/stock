package szse

import (
	"testing"
)

func TestSZSE(t *testing.T) {
	s := SZSE{Code: "002142"}
	if s.get().Realtime.Name != "宁波银行" {
		t.Error("Get szse stock error")
	}
}
