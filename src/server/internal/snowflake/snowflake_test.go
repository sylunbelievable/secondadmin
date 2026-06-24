package snowflake

import "testing"

func TestNextIncreases(t *testing.T) {
	g := New(7)
	first, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	second, err := g.Next()
	if err != nil {
		t.Fatal(err)
	}
	if second <= first {
		t.Fatalf("expected increasing ids, got %d then %d", first, second)
	}
}

func TestNextIsJavaScriptSafe(t *testing.T) {
	id, err := New(0).Next()
	if err != nil {
		t.Fatal(err)
	}
	if id > 1<<53-1 {
		t.Fatalf("id exceeds JavaScript safe integer: %d", id)
	}
}

func TestConfigureRejectsInvalidWorker(t *testing.T) {
	if err := Configure(1024); err == nil {
		t.Fatal("expected invalid worker id to fail")
	}
}
