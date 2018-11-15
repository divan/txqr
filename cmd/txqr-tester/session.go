package main

// SessionState defines test session state.
type SessionState int

const (
	StateNew SessionState = iota // default state
	StateStarted
	StatePaused
	StateFinished
	StateFailed
)

// Session represents a single test session.
type Session struct {
	state SessionState
}

// NewSession inits new test session.
// Default State is StateNew.
func NewSession() *Session {
	return &Session{}
}
