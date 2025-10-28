package logx_test

import (
	"os"
	"testing"
)

// TestMain resets the logs directory once per test run.
func TestMain(m *testing.M) {
	// Do not delete logs between runs; ensure directory exists for appending
	_ = os.MkdirAll("logs", 0755)
	os.Exit(m.Run())
}
