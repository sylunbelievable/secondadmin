package snowflake

import "testing"

func TestNextIncreases(t *testing.T) {
	if err := Configure(7); err != nil {
		t.Fatal(err)
	}
	first, err := Next()
	if err != nil {
		t.Fatal(err)
	}
	second, err := Next()
	if err != nil {
		t.Fatal(err)
	}
	if second <= first {
		t.Fatalf("expected increasing ids, got %d then %d", first, second)
	}
}

func TestConfigureRejectsInvalidWorker(t *testing.T) {
	if err := Configure(1024); err == nil {
		t.Fatal("expected invalid worker id to fail")
	}
}
