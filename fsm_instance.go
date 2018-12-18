package fsm

type FSM struct {
	machines map[string]*fsm
}

func NewFSM() *FSM {
	f := &FSM{}
	f.machines = make(map[string]*fsm)
	return f
}

func (f *FSM) Register(tag string, column string, events []EventTransition) error {
	f.machines[tag] = New(column, events)
	return nil
}

func (f *FSM) Set(tag string) *fsm {
	machine, ok := f.machines[tag]
	if !ok {
		return New("", Events{{}})
	}
	return machine
}
