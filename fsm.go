package fsm

import (
	"errors"
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
		f.guards[e.Name] = e.Guard
		for _, src := range e.From {

			f.transitions[eventKey{event: e.Name, src: src}] = e.To
		}
	}

	return f
}

func (f *FSM) Event(s interface{}, event string) error {

	val := reflect.ValueOf(s).Elem()

	if val.Kind() != reflect.Struct {
		return errors.New("struct not found")
	}

	v := val.FieldByName(f.column)

	if !v.CanSet() && v.Kind() != reflect.String {
		return errors.New("error")
	}

	destination, ok := f.transitions[eventKey{event, v.String()}]
	if !ok {
		return errors.New("event not found")
	}

	guard, ok := f.guards[event]

	ok = guard(s)
	if !ok {
		return errors.New("")
	}

	f.Lock()
	v.SetString(destination)
	f.Unlock()

	return nil
}
