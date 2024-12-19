package assert

import (
	"errors"
	"fmt"
	"testing"
)

func TestEqual(t *testing.T) {
	tst := &errorCapture{}
	Equal(tst, 1, 1)
	if tst.failed {
		t.Error("Equal failed")
	}

	tst = &errorCapture{}
	Equal(tst, 1, 2)
	if !tst.failed {
		t.Error("Equal failed")
	}
}

func TestNoError(t *testing.T) {
	tst := &errorCapture{}
	NoError(tst, nil)
	if tst.failed {
		t.Error("NoError failed")
	}

	tst = &errorCapture{}
	NoError(tst, errors.New("error"))
	if !tst.failed {
		t.Error("NoError failed")
	}
}

func TestError(t *testing.T) {
	tst := &errorCapture{}
	Error(tst, errors.New("error"), "error")
	if tst.failed {
		t.Error("Error failed")
	}

	tst = &errorCapture{}
	Error(tst, nil, "error")
	if !tst.failed {
		t.Error("Error failed")
	}

	tst = &errorCapture{}
	Error(tst, errors.New("error"), "other")
	if !tst.failed {
		t.Error("Error failed")
	}
}

func TestErrorIs(t *testing.T) {
	tst := &errorCapture{}
	ErrorIs(tst, errors.New("error"), errors.New("error"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	ErrorIs(tst, errors.New("error"), errors.New("other"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	ErrorIs(tst, nil, errors.New("error"))
	if !tst.failed {
		t.Error("ErrorIs failed")
	}

	tst = &errorCapture{}
	err := errors.New("error")
	ErrorIs(tst, fmt.Errorf("wrapped: %w", err), err)
	if tst.failed {
		t.Error("ErrorIs failed")
	}
}

func TestTrue(t *testing.T) {
	tst := &errorCapture{}
	True(tst, true)
	if tst.failed {
		t.Error("True failed")
	}

	tst = &errorCapture{}
	True(tst, false)
	if !tst.failed {
		t.Error("True failed")
	}
}

func TestFalse(t *testing.T) {
	tst := &errorCapture{}
	False(tst, false)
	if tst.failed {
		t.Error("False failed")
	}

	tst = &errorCapture{}
	False(tst, true)
	if !tst.failed {
		t.Error("False failed")
	}
}

func TestInterfaceNilEqual(t *testing.T) {
	tst := &errorCapture{}
	Equal(tst, nil, nil)
	if tst.failed {
		t.Error("InterfaceNilEqual failed")
	}

	tst = &errorCapture{}
	Equal(tst, nil, 1)
	if !tst.failed {
		t.Error("InterfaceNilEqual failed")
	}
}

func TestLen(t *testing.T) {
	tst := &errorCapture{}
	Len(tst, []int{1, 2}, 2)
	if tst.failed {
		t.Error("Len failed")
	}

	tst = &errorCapture{}
	Len(tst, []int{}, 2)
	if !tst.failed {
		t.Error("Len failed")
	}
}

func TestNotNil(t *testing.T) {
	tst := &errorCapture{}
	NotNil(tst, 1)
	if tst.failed {
		t.Error("NotNil failed")
	}

	tst = &errorCapture{}
	NotNil(tst, nil)
	if !tst.failed {
		t.Error("NotNil failed")
	}
}

func TestNil(t *testing.T) {
	tst := &errorCapture{}
	Nil(tst, nil)
	if tst.failed {
		t.Error("Nil failed")
	}

	tst = &errorCapture{}
	Nil(tst, 1)
	if !tst.failed {
		t.Error("Nil failed")
	}
}

func TestFail(t *testing.T) {
	tst := &errorCapture{}
	fail(tst, "error", "msg %d", 1)
	if !tst.failed {
		t.Error("Fail failed")
	}
	if tst.errs[0].(string) != "error\nmsg 1" {
		t.Error("Fail failed")
	}
}

type errorCapture struct {
	errs   []any
	failed bool
}

func (e *errorCapture) Helper() {
}

func (e *errorCapture) Error(args ...any) {
	e.errs = append([]any{}, args...)
}

func (e *errorCapture) FailNow() {
	e.failed = true
}
