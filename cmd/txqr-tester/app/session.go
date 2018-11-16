package main

import "github.com/divan/txqr/qr"

// SessionState defines test session state.
type SessionState int

// Predefined states
const (
	StateNew SessionState = iota // default state
	StateStarted
	StateAnimating
	StateWaitingNext
	StateFinished
	StateFailed
)

// Session represents a single test session.
type Session struct {
	state  SessionState
	config SessionConfig

	tests []*testSetup
	idx   int // current test index
}

// NewSession inits new test session.
// Default State is StateNew.
func NewSession() *Session {
	return &Session{
		config: DefaultSessionConfig(),
	}
}

// CurrentSetup returns the test setup of the currently executing round.
func (s *Session) CurrentSetup() testSetup {
	if s.state == StateNew || s.state == StateFinished {
		return testSetup{}
	}

	ts := s.tests[s.idx-1]
	return *ts
}

// UpdateConfig sets the new config for the session.
func (s *Session) UpdateConfig(config SessionConfig) {
	s.config = config
}

// StartNext starts next round of testing. It returns
// next untested parameters for QR code.
func (s *Session) StartNext() (*testSetup, bool) {
	if s.state == StateNew {
		s.state = StateStarted

		// generate parameters set, so they can't be changed during test
		s.tests = s.generateTestsSetup()
	}

	if s.idx == len(s.tests) {
		s.state = StateFinished
		return nil, false
	}

	test := s.tests[s.idx]
	s.idx++
	return test, true
}

// SetState explicitly sets the session state to state.
func (s *Session) SetState(state SessionState) {
	s.state = state
}

// State returns the current state of Session.
func (s *Session) State() SessionState {
	return s.state
}

// InProgress returns true if testing is in progress.
func (s *Session) InProgress() bool {
	return s.state != StateNew && s.state != StateFinished
}

func (s *Session) generateTestsSetup() []*testSetup {
	var ret []*testSetup
	for _, lvl := range AllQRLevels {
		if s.config.Levels[lvl] == false {
			continue
		}
		for fps := s.config.StartFPS; fps <= s.config.StopFPS; fps++ {
			for sz := s.config.StartSize; sz <= s.config.StopSize; sz += s.config.SizeStep {
				ret = append(ret, newTestSetup(fps, sz, lvl))
			}
		}
	}
	return ret
}

// testSetup represents setup parameters for the single test round.
type testSetup struct {
	fps  int
	size int
	lvl  qr.RecoveryLevel
}

func newTestSetup(fps, sz int, lvl qr.RecoveryLevel) *testSetup {
	return &testSetup{
		fps:  fps,
		size: sz,
		lvl:  lvl,
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

		Levels: DefaultQRLevels,
	}
}

// Levels is a handy wrapper type to work with a slice of RecoveryLevels.
type Levels map[qr.RecoveryLevel]bool

var DefaultQRLevels = map[qr.RecoveryLevel]bool{
	qr.Low:     true,
	qr.Medium:  true,
	qr.High:    true,
	qr.Highest: true,
}

var AllQRLevels = []qr.RecoveryLevel{
	qr.Low,
	qr.Medium,
	qr.High,
	qr.Highest,
}

func (levels Levels) has(lvl qr.RecoveryLevel) bool {
	return levels[lvl]
}

func (levels Levels) set(lvl qr.RecoveryLevel, val bool) {
	levels[lvl] = val
}

// numEnabled returns number of enabled levels.
func (levels Levels) numEnabled() int {
	var ret int
	for _, v := range levels {
		if v {
			ret++
		}
	}
	return ret
}
