[![Build Status](https://travis-ci.org/ssych/fsm.svg?branch=master)](https://travis-ci.org/ssych/fsm)

##  finite state machine for Go

This is package implements [finite state machine](https://en.wikipedia.org/wiki/Finite-state_machine)

# Basic Example

```go
package main

import (
  "fmt"
  "reflect"

  "github.com/ssych/fsm"
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

  f := fsm.NewFSM()

  f.Register(reflect.TypeOf((*Person)(nil)), "State", fsm.Events{{
    Name:   "grow_rich",
    From:   []fsm.State{"poor"},
    To:     fsm.State("rich"),
    Guard:  isRich,
    After:  after,
    Before: before,
  }})

  f.Fire(person, "grow_rich")
}

```
