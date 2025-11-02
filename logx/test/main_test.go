package logx_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = os.MkdirAll("logs", 0755)
	os.Exit(m.Run())
}
