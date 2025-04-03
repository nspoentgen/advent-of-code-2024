package main

import (
	"errors"
)

type Stack[TData any] struct {
	s []TData
}

func NewStack[TData any]() *Stack[TData] {
	return &Stack[TData]{make([]TData,0), }
}

func (s *Stack[TData]) Push(v TData) {
	s.s = append(s.s, v)
}

func (s *Stack[TData]) Pop() (TData, error) {
	l := len(s.s)
	if l == 0 {
		var dummy TData
		return dummy, errors.New("Empty Stack")
	}

	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func cloneMap[TKey comparable, TValue any](input map[TKey]TValue) map[TKey]TValue {
	clonedMap := make(map[TKey]TValue)

	for k, v := range input {
		clonedMap[k] = v
	}

	return clonedMap
}
