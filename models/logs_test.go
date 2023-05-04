package models

import (
	"errors"
	"testing"
)

func TestMockLogger_Create(t *testing.T) {
	ml := MockLogger{
		errorLog: []error{errors.New("first error"), errors.New("second error")},
	}

	newError := errors.New("third error to add")
	if err := ml.Create(newError); err != nil {
		t.Errorf("Mock logger returned error")
	}

	errorTable := ml.errorLog
	got := len(errorTable)
	want := 3
	if got != want {
		t.Errorf("The error was not appended to the error log, error log got size %d, want error log size %d",
			got, want,
		)

	}
	if errorTable[2].Error() != newError.Error() {
		t.Errorf("Order is wrong on error log")
	}
}
