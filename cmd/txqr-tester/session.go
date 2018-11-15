package main

import "github.com/divan/txqr/qr"

// SessionState defines test session state.
type SessionState int

// Predefined states
const (
	StateNew SessionState = iota // default state
	StateStarted
	StatePaused
	StateFinished
	StateFailed
)

// Session represents a single test session.
type Session struct {
	state  SessionState
	config SessionConfig
}

// NewSession inits new test session.
// Default State is StateNew.
func NewSession() *Session {
	return &Session{
		config: DefaultSessionConfig(),
	}
}

// SessionConfig holds configuration values for initiating the session.
type SessionConfig struct {
	StartFPS, StopFPS   int    // FPS of QR animation
	StartSize, StopSize int    // size of the chunk to be encoded into each animated QR frame
	SizeStep            int    // increment step for size changes
	Levels              Levels // recovery levels to use for test
}

// DefaultSessionConfig creates config with default values.
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		// from to 2 to 15 FPS
		StartFPS: 2,
		StopFPS:  15,

		// from 50 to 1000 with step 50
		StartSize: 50,
		StopSize:  1000,
		SizeStep:  50,

		Levels: AllQRLevels,
	}
}

// Levels is a handy wrapper type to work with a slice of RecoveryLevels.
type Levels []qr.RecoveryLevel

var AllQRLevels = Levels{
	qr.Low,
	qr.Medium,
	qr.High,
	qr.Highest,
}

func (levels Levels) has(lvl qr.RecoveryLevel) bool {
	for _, l := range levels {
		if l == lvl {
			return true
		}
	}
	return false
}
