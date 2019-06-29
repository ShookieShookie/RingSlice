package linked

import (
	"container/list"
	"errors"
)

// StaticList strict
type StaticList struct {
	*list.List
	start int
	next  int
	count int
}

// NewLinkedList new
func NewLinkedList(cap int) *StaticList {
	l := list.New()
	for i := 0; i < cap; i++ {
		l.PushBack(i + 1)
	}
	return &StaticList{List: l}
}

// Append new
func (s *StaticList) Append(val int) error {
	if s.count == s.Len() {
		return errors.New("Out of room")
	}
	current := s.GetIndex(s.trueIndex(s.start, s.count))
	current.Value = val
	s.count++
	return nil
}

func (s *StaticList) DeleteCount(length int) {
	if length > s.Len() {
		length = s.Len() // don't let us get stuck doing too much work
	}
	wipe := s.GetIndex(s.trueIndex(s.start, 0))
	for i := 0; i < length; i++ {
		wipe.Value = 0
		wipe = s.GetIndex(s.trueIndex(s.start, i+1))
	}
	s.start = s.trueIndex(s.start, length)
	s.count -= length
	if s.count < 0 { // TODO: could handle differntly, but perhaps valid to delete more than there is?
		s.count = 0
	}
}

func (s *StaticList) DeleteBounds(start, end int) {
	// get start node
	// wipe each node until end node
	// start node is next node after end node

	sn := s.GetIndex(s.trueIndex(start, 0))
	en := s.GetIndex(s.trueIndex(end, 0))

	for cur := sn; ; cur = cur.Next() {
		if cur == nil {
			cur = s.Front()
		}
		cur.Value = 0
		if s.count > 0 {
			s.count--
		}
		if cur == en {
			break
		}
	}

	s.start = s.trueIndex(end, 1) // one past the last thing we deleted

}

func (s *StaticList) trueIndex(start, length int) int {
	return (start + length) % s.Len()
}

func (s *StaticList) GetIndex(ind int) *list.Element {
	n := s.Front()
	for i := 0; i < ind; i++ {
		n = n.Next()
	}
	return n
}

func (s *StaticList) Values() []int {
	vals := []int{}
	node := s.Front()
	for i := 0; i < s.Len(); i++ {
		vals = append(vals, node.Value.(int))
		node = node.Next()
	}
	return vals
}
