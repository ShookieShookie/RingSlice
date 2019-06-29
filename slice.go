package main

import (
	"errors"
	"fmt"
)

// Slice struct
type Slice struct {
	values []int
	used   int
	start  int
	end    int
	debug  bool
	cap    int
}

// TODO provide a clear function

// NewSlice does
func NewSlice(capacity int, debug bool) *Slice {
	return &Slice{values: make([]int, capacity), debug: debug, cap: capacity}
}

// Append does
func (s *Slice) Append(value int) error {
	if s.used == s.cap {
		return errors.New("Index is full cannot append")
	}
	ind := s.trueIndex(s.start, s.used) // next index is same as num written
	if s.debug {
		fmt.Println("ind to write to ", ind)
	}
	s.values[ind] = value
	s.used++
	return nil
}

// Fetch does
func (s *Slice) Fetch(index int) int {
	return s.values[index]
}

// DeleteBounds does
func (s *Slice) DeleteBounds(start, end int) {
	stop := s.trueIndex(end, 0)
	for i := start; ; i = s.next(i) {
		s.values[i] = 0
		s.used-- // TODO: could speed this up with calculation
		if i == stop {
			break
		}
	}
	s.start = stop
}

// DeleteCount does
func (s *Slice) DeleteCount(count int) {
	ind := s.start
	if count > s.used {
		count = s.used // save us some time
	}
	for i := 0; i < count; i++ {
		s.values[ind] = 0
		ind = s.next(ind)
	}
	s.used -= count
	s.start = ind
}

func (s *Slice) next(cur int) int {
	if cur == s.cap-1 {
		return 0
	}
	return cur + 1
}

// returns index of length AWAY from start taking into account wrap around
// start,0 == start
func (s *Slice) trueIndex(start, length int) int {
	return (start + length) % s.cap
}

func (s *Slice) firstFree() int {
	for i := 0; i < len(s.values); i++ {
		if s.values[i] == 0 { // TODO obviously this is a poor check, need an allocation variable?
			return i
		}
	}
	return -1
}
