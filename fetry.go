package fetry

import (
	"errors"
	"reflect"
	"time"
)

var (
	ErrExecNotFunction = errors.New("exec type not function")
	ErrArgsLenNotMatch = errors.New("input arguments length not match")
	ErrOutsLenNotMatch = errors.New("output arguments length not match")
	ErrOutsTypeNotPtrs = errors.New("output arguments type not pointer")
	ErrOutputArgNotNil = errors.New("output arguments not nil")
)

type Fetry struct {
	exec     interface{}
	args     []interface{}
	times    uint
	interval time.Duration
}

func NewFetry(f interface{}, times uint, interval time.Duration, args ...interface{}) *Fetry {
	r := new(Fetry)
	r.exec = f
	r.args = args
	r.times = times
	r.interval = interval
	return r
}

func (f *Fetry) Exec() error {
	exec := reflect.ValueOf(f.exec)
	t := exec.Type()

	// check if exec is function
	if t.Kind() != reflect.Func {
		return ErrExecNotFunction
	}

	// check input arguments length
	if t.NumIn() > len(f.args) {
		return ErrArgsLenNotMatch
	}

	// check output arguments length
	if t.NumOut() != 1 {
		return ErrOutsLenNotMatch
	}

	// handle input arguments
	length := len(f.args)
	inputs := make([]reflect.Value, length)
	for i := 0; i < length; i++ {
		if f.args[i] == nil {
			inputs[i] = reflect.Zero(exec.Type().In(i))
		} else {
			inputs[i] = reflect.ValueOf(f.args[i])
		}
	}

	// exec function
	output := exec.Call(inputs)

	// handle output arguments
	if len(output) != 1 {
		return ErrOutsLenNotMatch
	}
	if output[0].IsNil() {
		return nil
	}
	return ErrOutputArgNotNil
}
