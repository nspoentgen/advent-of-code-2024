package main

import (
	"errors"
	"sync"
)

type Stack[TData any] struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []TData
}

func NewStack[TData any]() *Stack[TData] {
	return &Stack[TData]{sync.Mutex{}, make([]TData, 0)}
}

func (s *Stack[TData]) Push(v TData) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *Stack[TData]) Pop() (TData, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		var dummy TData
		return dummy, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *Stack[TData]) Len() int {
	return len(s.s)
}
