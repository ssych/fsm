package fsm

import (
	"reflect"
	"sync"
)

type Guard func(*Event) (bool, error)

type Event struct {
	Event       string
	Source      interface{}
	Destination State
}

type EventTransition struct {
	Name   string
	From   []State
	To     State
	Guards []Guard
	After  func(*Event) error
	Before func(*Event) error
}

type Events []EventTransition

type fsm struct {
	sync.RWMutex
	column        string
	transitions   map[eventKey]State
	initialStates map[State][]string
	guards        map[string][]Guard
	callbacks     map[cKey]func(*Event) error
}

type eventKey struct {
	event string
	src   State
}

type cKey struct {
	event string
	cType string
}

func newFSM(column string, events []EventTransition) *fsm {
	f := &fsm{
		column: column,
	}
	f.transitions = make(map[eventKey]State)
	f.guards = make(map[string][]Guard)
	f.callbacks = make(map[cKey]func(*Event) error)
	f.initialStates = make(map[State][]string)

	for _, e := range events {
		if e.Guards != nil {
			f.guards[e.Name] = e.Guards
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

	for eventKey := range f.transitions {
		f.initialStates[eventKey.src] = append(f.initialStates[eventKey.src], eventKey.event)
	}

	return f
}

func (f *fsm) Fire(s interface{}, event string) error {
	state, err := f.getSourceState(s)
	if err != nil {
		return err
	}

	destination, ok := f.transitions[eventKey{event, State(state.String())}]
	if !ok {
		return UnknownEventError{event}
	}

	e := &Event{Event: event, Source: s, Destination: destination}

	ok, err = f.guardEvent(e)
	if err != nil {
		return err
	}

	if !ok {
		return InvalidTransitionError{event, state.String()}
	}

	f.Lock()

	err = f.beforeEventCallbacks(e)
	if err != nil {
		return err
	}

	state.SetString(string(destination))

	err = f.afterEventCallbacks(e)
	if err != nil {
		return err
	}

	f.Unlock()

	return nil
}

func (f *fsm) MayFire(s interface{}, event string, options ...Option) (bool, error) {
	// Setup options.
	args := &Options{}
	for _, option := range options {
		option(args)
	}

	state, err := f.getSourceState(s)
	if err != nil {
		return false, err
	}

	destination, ok := f.transitions[eventKey{event, State(state.String())}]
	if !ok {
		return false, nil
	}

	e := &Event{Event: event, Source: s, Destination: destination}

	if !args.SkipGuards {
		ok, err = f.guardEvent(e)
		if err != nil {
			return false, err
		}
	}

	return ok, nil
}

func (f *fsm) GetPermittedEvents(s interface{}, options ...Option) ([]string, error) {
	state, err := f.getSourceState(s)
	if err != nil {
		return nil, err
	}

	events, ok := f.initialStates[State(state.String())]
	if !ok {
		return []string{}, nil
	}

	permittedEvents := []string{}
	for _, event := range events {
		ok, err := f.MayFire(s, event, options...)
		if err != nil {
			return nil, err
		}

		if ok {
			permittedEvents = append(permittedEvents, event)
		}
	}

	return permittedEvents, nil
}

func (f *fsm) getSourceState(s interface{}) (state reflect.Value, err error) {
	val := reflect.ValueOf(s).Elem()

	if val.Kind() != reflect.Struct {
		return state, InternalError{}
	}

	state = val.FieldByName(f.column)
	if !state.IsValid() && !state.CanSet() && state.Kind() != reflect.String {
		return state, InternalError{}
	}

	return
}

func (f *fsm) guardEvent(e *Event) (bool, error) {
	fns, ok := f.guards[e.Event]
	if ok {
		for _, fn := range fns {
			if ok, err := fn(e); err != nil || !ok {
				return false, err
			}
		}
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
