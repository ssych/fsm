package fsm

import (
	"reflect"
	"sync"
)

type Event struct {
	Name  string
	From  []string
	To    string
	Guard func(interface{}) bool
}

type Events []Event

type FSM struct {
	sync.RWMutex
	column      string
	transitions map[eventKey]string
	guards      map[string]func(interface{}) bool
}

type eventKey struct {
	event string
	src   string
}

func New(column string, events []Event) *FSM {
	f := &FSM{
		column: column,
	}
	f.transitions = make(map[eventKey]string)
	f.guards = make(map[string]func(interface{}) bool)

	for _, e := range events {
		if e.Guard != nil {
			f.guards[e.Name] = e.Guard
		}

		for _, src := range e.From {
			f.transitions[eventKey{event: e.Name, src: src}] = e.To
		}
	}

	return f
}

func (f *FSM) Fire(s interface{}, event string) error {

	val := reflect.ValueOf(s).Elem()

	if val.Kind() != reflect.Struct {
		return InternalError{}
	}

	state := val.FieldByName(f.column)

	if !state.IsValid() && !state.CanSet() && state.Kind() != reflect.String {
		return InternalError{}
	}

	destination, ok := f.transitions[eventKey{event, state.String()}]
	if !ok {
		return UnknownEventError{event}
	}

	guard, ok := f.guards[event]
	if ok && !guard(s) {
		return InvalidTransitionError{event, state.String()}
	}

	f.Lock()
	state.SetString(destination)
	f.Unlock()

	return nil
}
