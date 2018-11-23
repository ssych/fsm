package fsm

import (
	"testing"
)

type TestStruct struct {
	State State
}

func IsTestStructValid(s interface{}) bool {
	return true
}

func IsTestStructInvalid(s interface{}) bool {
	return false
}

func TestSetState(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := New("State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructValid,
	}})

	fsm.Event(testStruct, "make")
	if testStruct.State != State("finished") {
		t.Error("expected state to be 'finished'")
	}
}

func TestInvalidTransition(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := New("State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Event(testStruct, "make")

	if e, ok := err.(InvalidTransitionError); !ok && e.Event != "make" && e.State != "started" {
		t.Error("expected 'InvalidTransitionError'")
	}
}

func TestInvalidEvent(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := New("State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Event(testStruct, "some_event_name")

	if e, ok := err.(UnknownEventError); !ok && e.Event != "some_event_name" {
		t.Error("expected 'UnknownEventError'")
	}
}
