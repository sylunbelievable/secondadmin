package dto

import (
	"encoding/json"
	"testing"
)

func TestIDsRequestAcceptsStringAndNumberIDs(t *testing.T) {
	tests := []string{
		`{"ids":["1999999999999999999","2"]}`,
		`{"ids":[1999999999999999999,2]}`,
	}
	for _, input := range tests {
		var req IDsRequest
		if err := json.Unmarshal([]byte(input), &req); err != nil {
			t.Fatalf("unmarshal %s: %v", input, err)
		}
		if got := []uint64(req.IDs); got[0] != 1999999999999999999 || got[1] != 2 {
			t.Fatalf("unexpected ids: %#v", got)
		}
	}
}
