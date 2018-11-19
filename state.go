package fsm

import (
	"database/sql/driver"
	"log"
)

type State string

func (s *State) Scan(value interface{}) error {
	var str string
	switch t := value.(type) {
	case []uint8:
		str = string([]byte(value.([]uint8)))
	case string:
		str = value.(string)
	default:
		log.Fatalf("Scan(): unexpected type: %T for %#v", t, value)
	}
	*s = State(str)
	return nil
}

func (s *State) Value() (driver.Value, error) {
	return s, nil
}
