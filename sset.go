package fetry

import (
	"sort"
	"sync"
)

type Score []int64

func (s Score) Len() int { return len(s) }

func (s Score) Less(i, j int) bool { return s[i] < s[j] }

func (s Score) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *Score) Del(i int) {
	if i < 0 || len(*s) <= i {
		return
	}
	*s = append((*s)[:i], (*s)[i+1:]...)
}

func (s *Score) Search(x int64) int {
	left := 0
	right := s.Len() - 1
	middle := (left + right) >> 1
	for left < right {
		if (*s)[middle] == x {
			return middle
		} else if (*s)[middle] > x {
			right = middle - 1
		} else {
			left = middle + 1
		}
		middle = (left + right) >> 1
	}
	return -1
}

type SortedSet struct {
	sync.RWMutex
	scores Score
	values map[int64]interface{}
}

func NewSortedSet() *SortedSet {
	var s SortedSet
	s.values = make(map[int64]interface{})
	return &s
}

func (ss *SortedSet) Push(score int64, value interface{}) {
	ss.Lock()
	defer ss.Unlock()
	ss.scores = append(ss.scores, score)
	ss.values[score] = value
	sort.Sort(ss.scores)
}

func (ss *SortedSet) Pop() (int64, interface{}, bool) {
	ss.Lock()
	defer ss.Unlock()
	if len(ss.scores) == 0 {
		return 0, nil, false
	}
	s := ss.scores[0]
	v, b := ss.values[s]
	ss.scores = ss.scores[1:]
	return s, v, b
}

func (ss *SortedSet) DelByIndex(index int) {
	ss.Lock()
	defer ss.Unlock()
	if len(ss.scores) <= index {
		return
	}
	delete(ss.values, ss.scores[index])
	ss.scores.Del(index)
}

func (ss *SortedSet) DelByScore(score int64) {
	ss.Lock()
	defer ss.Unlock()
	delete(ss.values, score)
	ss.scores.Del(ss.scores.Search(score))
}

func (ss *SortedSet) GetByIndex(index int) (interface{}, bool) {
	ss.RLock()
	defer ss.RUnlock()
	if len(ss.scores) <= index {
		return nil, false
	}
	v, b := ss.values[ss.scores[index]]
	return v, b
}

func (ss *SortedSet) GetByScore(score int64) (interface{}, bool) {
	ss.RLock()
	defer ss.RUnlock()
	v, b := ss.values[score]
	return v, b
}

func (ss *SortedSet) Len() int {
	ss.RLock()
	defer ss.RUnlock()
	return len(ss.scores)
}

func (ss *SortedSet) Empty() bool {
	ss.RLock()
	defer ss.RUnlock()
	return len(ss.scores) == 0
}
