# Fetry

A golang function retry library.

## Install

* go get github.com/wangh10/fetry

### Usage

There are two methods which need retry: f1, f2. f1 retry 9 times. f2 retry 3 times.
```go
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
```
