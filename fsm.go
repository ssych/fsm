package fsm

import (
	"reflect"
	"sync"
)

type Event struct {
	Event       string
	Source      interface{}
	Destination string
}

type EventTransition struct {
	Name   string
	From   []string
	To     string
	Guard  func(*Event) (bool, error)
	After  func(*Event) error
	Before func(*Event) error
}

type Events []EventTransition

type fsm struct {
	sync.RWMutex
	column      string
	transitions map[eventKey]string
	guards      map[string]func(*Event) (bool, error)
	callbacks   map[cKey]func(*Event) error
}

type eventKey struct {
	event string
	src   string
}

type cKey struct {
	event string
	cType string
}

func New(column string, events []EventTransition) *fsm {
	f := &fsm{
		column: column,
	}
	f.transitions = make(map[eventKey]string)
	f.guards = make(map[string]func(*Event) (bool, error))
	f.callbacks = make(map[cKey]func(*Event) error)

	for _, e := range events {
		if e.Guard != nil {
			f.guards[e.Name] = e.Guard
		}

		if e.After != nil {
			f.callbacks[cKey{event: e.Name, cType: "after"}] = e.After
		}

		if e.Before != nil {
			f.callbacks[cKey{event: e.Name, cType: "before"}] = e.Before
		}

		for _, src := range e.From {
			f.transitions[eventKey{event: e.Name, src: src}] = e.To
		}
	}

	return f
}

func (f *fsm) Fire(s interface{}, event string) error {

	val := reflect.ValueOf(s).Elem()

	if val.Kind() != reflect.Struct {
		return InternalError{}
	}

	state := val.FieldByName(f.column)

	if !state.IsValid() && !state.CanSet() && state.Kind() != reflect.String {
		return InternalError{}
	}

	src := state.String()

	destination, ok := f.transitions[eventKey{event, src}]
	if !ok {
		return UnknownEventError{event}
	}

	e := &Event{Event: event, Source: s, Destination: destination}

	ok, err := f.guardEvent(e)

	if err != nil {
		return err
	}

	if !ok {
		return InvalidTransitionError{event, src}
	}

	f.Lock()

	err = f.beforeEventCallbacks(e)
	if err != nil {
		return err
	}

	state.SetString(destination)

	err = f.afterEventCallbacks(e)
	if err != nil {
		return err
	}

	f.Unlock()

	return nil
}

func (f *fsm) guardEvent(e *Event) (bool, error) {
	fn, ok := f.guards[e.Event]
	if ok {
		return fn(e)
	}
	return true, nil
}

func (f *fsm) afterEventCallbacks(e *Event) error {
	fn, ok := f.callbacks[cKey{event: e.Event, cType: "after"}]
	if ok {
		return fn(e)
	}
	return nil
}

func (f *fsm) beforeEventCallbacks(e *Event) error {
	fn, ok := f.callbacks[cKey{event: e.Event, cType: "before"}]
	if ok {
		return fn(e)
	}
	return nil
}
