package capitalflows

import "testing"

func TestFlows(t *testing.T) {
	flows, err := Fetch()
	if err != nil {
		t.Skip(err)
		return
	}
	t.Log(len(flows), flows)
}
