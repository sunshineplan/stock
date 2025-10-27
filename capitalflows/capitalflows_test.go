package capitalflows

import (
	"log"
	"testing"
)

func TestFlows(t *testing.T) {
	flows, err := Fetch()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(len(flows), flows)
}
