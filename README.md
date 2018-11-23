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

func IsRich(s interface{}) bool {
  person := s.(*Person)
  if person.Сash > 1000000 {
    return true
  }
  return false
}

func main() {

  person := &Person{
    Сash:  1000001,
    State: fsm.State("poor"),
  }

  fsm := fsm.New("State", fsm.Events{{
    Name:  "grow_rich",
    From:  []string{"poor"},
    To:    "rich",
    Guard: IsRich,
  }})

  fmt.Println(person.State)
  fsm.Fire(person, "grow_rich")
  fmt.Println(person.State)
}
```