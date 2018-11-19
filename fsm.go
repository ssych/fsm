package fsm

import (
	"errors"
	"reflect"
	"sync"
)

type Event struct {
	Name string
	From []string
	To   string
}

type Events []Event

type FSM struct {
	sync.RWMutex
	column      string
	transitions map[eventKey]string
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

	for _, e := range events {
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
		return errors.New("error types")
	}

	destination, ok := f.transitions[eventKey{event, v.String()}]
	if !ok {
		return errors.New("event not found")
	}

	f.Lock()
	v.SetString(destination)
	f.Unlock()

	return nil
}
