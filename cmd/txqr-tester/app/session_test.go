package main

import "testing"

func TestGenerateTests(t *testing.T) {
	sess := NewSession()
	got := sess.generateTestsSetup()
	expectedLen := 988
	if len(got) != expectedLen {
		t.Fatalf("Expected len to be %d, but got %d", expectedLen, got)
	}
}
