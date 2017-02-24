package fetry

import "time"

type Queue interface {
	Push(f *Fetry)
	Exit()
}

type queue struct {
	ss *SortedSet
	cc chan struct{}
}

func NewQueue() Queue {
	var q queue
	q.cc = make(chan struct{}, 1)
	q.ss = NewSortedSet()
	go q.exec()
	return &q
}

func (q *queue) Push(f *Fetry) {
	q.ss.Push(time.Now().UnixNano(), f)
}

func (q *queue) Exit() {
	q.cc <- struct{}{}
}

func (q *queue) exec() {
	for {
		select {
		case <-q.cc:
			return
		default:
			srt, pop, has := q.ss.Pop()
			if !has {
				continue
			}
			fun, ok := pop.(*Fetry)
			if !ok {
				continue
			}
			if time.Now().UnixNano() < srt {
				q.ss.Push(srt, fun)
				continue
			}
			if fun.times == 0 {
				continue
			}
			go func(f *Fetry) {
				e := f.Exec()
				if e == ErrOutputArgNotNil {
					f.times--
					q.ss.Push(time.Now().UnixNano()+f.interval.Nanoseconds(), f)
				}
			}(fun)
		}
	}
}
