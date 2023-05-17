package test

import (
	"io"
	"time"
)

// A TestInfo contains various information about a test.
type TestInfo struct {
	data           interface{}
	testID         string
	timeout        time.Duration
	header         io.Writer
	writer         io.Writer
	suite          *suiteInfo
}

// SetupInfo returns the eventual object stored by the Setup function.
func (t TestInfo) SetupInfo() interface{} {
	return t.data
}

// SuiteSetupInfo returns the eventual object stored by the Suite Setup function.
func (t TestInfo) SuiteSetupInfo() interface{} {
	return t.suite.data
}

// Stash provides the callers data that may store configs/runtime they need
func (t TestInfo) Stash() interface{} {
	return t.suite.stash
}

// Iteration returns the test iteration number.
func (t TestInfo) Iteration() int {
	return t.iteration
}

// TestID returns the test ID
func (t TestInfo) TestID() string {
	return t.testID
}

// WriteHeader performs a write at the header
func (t TestInfo) WriteHeader(p []byte) (n int, err error) {
	return t.header.Write(p)
}

// Write performs a write
func (t TestInfo) Write(p []byte) (n int, err error) {
	return t.writer.Write(p)
}

// TimeSinceLastStep provides the time since last step or assertion
func (t TestInfo) TimeSinceLastStep() string {
	d := time.Since(t.timeOfLastStep)
	return d.Round(time.Millisecond).String()
}

// Timeout provides the duration before the test timeout.
func (t TestInfo) Timeout() time.Duration {
	return t.timeout
}
