package fsm

type InvalidTransitionError struct {
	Event string
	State string
}

func (e InvalidTransitionError) Error() string {
	return "Event " + e.Event + "cannot transition from " + e.State
}

type UnknownEventError struct {
	Event string
}

func (e UnknownEventError) Error() string {
	return "event " + e.Event + " does not exist"
}

type InternalError struct{}

func (InternalError) Error() string {
	return "internal error"
}
