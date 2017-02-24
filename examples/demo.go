package main

import (
	"errors"
	"time"

	"github.com/wangh10/fetry"
)

func Div(a, b int) error {
	println(a, "/", b)
	if b == 0 {
		return errors.New("div zero")
	}
	return nil
}

func main() {
	f1 := fetry.NewFetry(Div, 9, time.Second, 1, 0)
	f2 := fetry.NewFetry(Div, 3, 3*time.Second, 3, 0)
	q := fetry.NewQueue()
	q.Push(f1)
	time.Sleep(2 * time.Second)
	q.Push(f2)
	time.Sleep(time.Minute)
	q.Exit()
}
