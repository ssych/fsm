##  finite state machine for Go

This is package implements [finite state machine](https://en.wikipedia.org/wiki/Finite-state_machine)

# Basic Example

```go
package main

import (
  "fmt"

  "github.com/stasmoon/fsm"
)

type Person struct {
  Сash  int
  State fsm.State
}

func isRich(e *fsm.Event) (bool, error) {
  person := e.Source.(*Person)
  if person.Сash > 1000000 {
    return true, nil
  }
  return false, nil
}

func after(e *fsm.Event) error {
  person := e.Source.(*Person)
  fmt.Println(person.State)
  return nil
}

func before(e *fsm.Event) error {
  person := e.Source.(*Person)
  fmt.Println(person.State)
  return nil
}

func main() {

  person := &Person{
    Сash:  1000001,
    State: fsm.State("poor"),
  }

  fsm := fsm.New("State", fsm.Events{{
    Name:   "grow_rich",
    From:   []string{"poor"},
    To:     "rich",
    Guard:  isRich,
    After:  after,
    Before: before,
  }})

  fsm.Fire(person, "grow_rich")
}

```