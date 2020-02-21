package fsm

import (
	"reflect"
)

type FSM struct {
	machines map[reflect.Type]*fsm
}

// NewFSM func to create FSM
func NewFSM() *FSM {
	f := &FSM{}
	f.machines = make(map[reflect.Type]*fsm)
	return f
}

// Register func to register all event by model reflect type
func (f *FSM) Register(tag reflect.Type, column string, events []EventTransition) error {
	f.machines[tag] = newFSM(column, events)
	return nil
}

// Fire func to fire event
func (f *FSM) Fire(s interface{}, event string) error {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return InternalError{}
	}

	return machine.Fire(s, event)
}

// MayFire func return false if event can`t may fire
func (f *FSM) MayFire(s interface{}, event string, options ...Option) (bool, error) {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return false, InternalError{}
	}

	return machine.MayFire(s, event, options...)
}

// GetPermittedEvents func to return all permitted events
func (f *FSM) GetPermittedEvents(s interface{}, options ...Option) ([]string, error) {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return nil, InternalError{}
	}

	return machine.GetPermittedEvents(s, options...)
}

// GetPermittedStates func to return all permitted states
func (f *FSM) GetPermittedStates(s interface{}, options ...Option) ([]State, error) {
	machine, ok := f.machines[reflect.TypeOf(s)]
	if !ok {
		return nil, InternalError{}
	}

	return machine.GetPermittedStates(s, options...)
}
