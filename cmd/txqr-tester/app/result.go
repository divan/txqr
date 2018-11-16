package main

import "time"

// Result represents test result from client.
type Result struct {
	testSetup
	Duration time.Duration
}

// NewResult constructs a new result.
func NewResult(setup testSetup, d time.Duration) Result {
	return Result{
		testSetup: setup,
		Duration:  d,
	}
}
