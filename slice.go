package ringslice

import (
	"errors"
	"fmt"
)

// Slice struct
type Slice struct {
	values []interface{}
	used   int
	start  int
	cap    int
	wipe   func(int, []interface{})
}

// NewSlice does
func NewSlice(capacity int, debug bool, wipe func(int, []interface{})) *Slice {
	return &Slice{values: make([]interface{}, capacity), cap: capacity, wipe: wipe}
}

// Append adds an entry if possible, returns error if full
func (s *Slice) Append(value interface{}) error {
	if s.used == s.cap {
		return errors.New("Index is full cannot append")
	}
	ind := s.trueIndex(s.start, s.used) // next index is same as num written
	s.values[ind] = value
	s.used++
	return nil
}

// Values provides a set of values for debugging purposes by taking ring
// and applying valuation function to each entry
func (s *Slice) Values(value func(interface{}) int64) []int64 {
	v := []int64{}
	for _, val := range s.values {
		v = append(v, value(val))
	}
	return v
}

// Stats prints information for debugging the slice
func (s *Slice) Stats(value func(interface{}) int64) {
	fmt.Println("used", s.used, "start", s.start, "valid", s.validate(value))
}

func (s *Slice) validate(value func(interface{}) int64) error {
	if s.start > (s.cap - 1) {
		return errors.New("illegal start")
	}
	if s.used <= 1 {
		return nil
	}
	var last int64
	for i := s.start; i < s.start+s.used-1; i++ {
		next := value(s.values[s.trueIndex(i, 0)])
		if next < last {
			return errors.New("values out of order")
		}
	}
	return nil
}

// Purge wipes all indices that have a value determined by value function
// to be <= want
// TODO keep track of min and max whether we should even check
func (s *Slice) Purge(want int64, value func(interface{}) int64) {
	ind := s.FindClosestBelowOrEqual(want, value)
	if ind == -1 {
		return
	}
	fmt.Println("deleting bounds", s.start, ind)
	s.DeleteBounds(s.start, ind)
}

// FindClosestBelowOrEqual uses a binary search to find the HIGHEST value that is <= want
// if nothing is <= value, return -1. Accounts for array wrapping around by determining the bounds
// of the array if it were laid out contiguously
func (s *Slice) FindClosestBelowOrEqual(want int64, value func(interface{}) int64) int {
	if s.used == 0 {
		return -1
	}
	// binary search and call value on node to find value
	start := s.trueIndex(s.start, 0)
	falseMax := start + s.used - 1 // if the slice were continous, the highest index
	end := falseMax

	if end == start { // there's one node
		if value(s.values[s.trueIndex(start, 0)]) > want {
			return -1 // don't delete it, it's too low
		}
		return start
	}
	// start at midpoint, always lookup with index value % length of array
	for m := (start + falseMax) / 2; ; {
		cur := value(s.values[s.trueIndex(m, 0)])
		if cur == want {
			// if we find the exact value, walk to latest index with this value
			return s.findLatestEquivalent(m, want, value)
		}

		// if the start and end indices are 1 away, return higher index <= want
		if end-start == 1 || (end == 0 && start == s.cap-1) {
			return s.determineBoundary(start, end, want, value)
		}

		// if the value we check was less, set the start of next midpoint here
		if cur < want {
			start = m
		}
		// if the value we check was greater, set end to be this point
		if cur > want {
			end = m
		}
		// set the next check index to be midpoint between changed start and end
		m = (start + end) / 2
	}
}

// findLatestEquivalent walks clockwise in the array until the value changes, finding the highest
// value <= want linearly. TODO: Could be optimized (val-count map) but intended case does not have any/ few equals.
func (s *Slice) findLatestEquivalent(m int, want int64, value func(interface{}) int64) int {
	for new := m; ; {
		new = s.next(new)
		if new == m {
			return s.trueIndex(s.start+s.used-1, 0) // we've iterated back to where we started, all values equal
		}
		if value(s.values[s.trueIndex(new, 0)]) != want {
			return s.prev(new)
		}
	}
}

// determineBoundary provides logic for case when pointers are 1 away from each other
// happens when 1) all values are > 2) all values are < 3) there is a set below and a set above want
func (s *Slice) determineBoundary(start, end int, want int64, value func(interface{}) int64) int {
	if value(s.values[s.trueIndex(end, 0)]) <= want {
		return end
	}
	if value(s.values[s.trueIndex(start, 0)]) > want {
		return -1
	}
	return start
}

// DeleteBounds deletes all indices [start,end] inclusive
func (s *Slice) DeleteBounds(start, end int) {
	stop := s.trueIndex(end, 0)
	for i := start; ; i = s.next(i) {
		s.wipe(i, s.values)
		if s.used > 0 {
			s.used-- // TODO: could speed this up with calculation
		}
		if i == stop {
			break
		}
	}
	s.start = s.next(stop)
}

// DeleteCount deletes count of values starting at start index
func (s *Slice) DeleteCount(count int) {
	ind := s.start
	if count > s.used {
		count = s.used // save us some time
	}
	for i := 0; i < count; i++ {
		s.wipe(i, s.values)
		ind = s.next(ind)
	}
	s.used -= count
	s.start = ind
}

// safely iterate loop clockwise
func (s *Slice) next(cur int) int {
	if cur == s.cap-1 {
		return 0
	}
	return cur + 1
}

// safely iterate loop counterclockwise
func (s *Slice) prev(cur int) int {
	if cur == 0 {
		return s.cap - 1
	}
	return cur - 1
}

// returns index of length AWAY from start taking into account wrap around
// start,0 == start
func (s *Slice) trueIndex(start, length int) int {
	return (start + length) % s.cap
}
