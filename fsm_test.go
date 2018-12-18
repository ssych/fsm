package fsm

import (
	"testing"
)

type TestStruct struct {
	State State
}

func IsTestStructValid(e *Event) (bool, error) {
	return true, nil
}

func IsTestStructInvalid(e *Event) (bool, error) {
	return false, nil
}

func TestSetState(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()

	fsm.Register("test", "State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructValid,
	}})

	fsm.Set("test").Fire(testStruct, "make")
	if testStruct.State != State("finished") {
		t.Error("expected state to be 'finished'")
	}
}

func TestInvalidTransition(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()

	fsm.Register("test", "State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Set("test").Fire(testStruct, "make")

	if e, ok := err.(InvalidTransitionError); !ok && e.Event != "make" && e.State != "started" {
		t.Error("expected 'InvalidTransitionError'")
	}
}

func TestInvalidEvent(t *testing.T) {
	testStruct := &TestStruct{
		State: State("started"),
	}

	fsm := NewFSM()
	fsm.Register("test", "State", Events{{
		Name:  "make",
		From:  []string{"started"},
		To:    "finished",
		Guard: IsTestStructInvalid,
	}})

	err := fsm.Set("test").Fire(testStruct, "some_event_name")

	if e, ok := err.(UnknownEventError); !ok && e.Event != "some_event_name" {
		t.Error("expected 'UnknownEventError'")
	}
}
