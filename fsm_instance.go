package fsm

import (
	"reflect"
)

type FSM struct {
	machines map[reflect.Type]*fsm
}

func NewFSM() *FSM {
	f := &FSM{}
	f.machines = make(map[reflect.Type]*fsm)
	return f
}

func (f *FSM) Register(tag reflect.Type, column string, events []EventTransition) error {
	f.machines[tag] = newFSM(column, events)
	return nil
}

func (f *FSM) Fire(s interface{}, event string) error {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return InternalError{}
	}

	return machine.Fire(s, event)
}

func (f *FSM) MayFire(s interface{}, event string, options ...Option) (bool, error) {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return false, InternalError{}
	}

	return machine.MayFire(s, event, options...)
}

func (f *FSM) GetPermittedEvents(s interface{}, options ...Option) ([]string, error) {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return nil, InternalError{}
	}

	return machine.GetPermittedEvents(s, options...)
}
