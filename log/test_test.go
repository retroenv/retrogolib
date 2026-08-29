package log

import (
	"fmt"
	"strings"
	"testing"

	"github.com/retroenv/retrogolib/assert"
)

func TestNewTestLogger(t *testing.T) {
	logger := NewTestLogger(t)
	assert.Equal(t, TraceLevel, logger.level.Level())
}

func TestTestLoggerChildrenPreserveFailureBehavior(t *testing.T) {
	test := &recordingTest{}
	logger := NewTestLogger(test).With(String("component", "cpu")).Named("step")

	logger.Trace("Trace")
	logger.Error("Failed")

	assert.True(t, test.failed)
	assert.Contains(t, strings.Join(test.logs, "\n"), "Trace")
	assert.Contains(t, strings.Join(test.logs, "\n"), "Failed")
}

type recordingTest struct {
	failed bool
	logs   []string
}

func (t *recordingTest) Logf(format string, args ...any) {
	t.logs = append(t.logs, fmt.Sprintf(format, args...))
}

func (t *recordingTest) FailNow() {
	t.failed = true
}

func (t *recordingTest) Helper() {}
