package main

import "testing"

func TestAdminMenuSeeds(t *testing.T) {
	if err := validateAdminMenuSeeds(); err != nil {
		t.Fatal(err)
	}
	if len(adminMenuSeeds) != 28 {
		t.Fatalf("expected 28 admin menu nodes, got %d", len(adminMenuSeeds))
	}
}
