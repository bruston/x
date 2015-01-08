package intset

import (
	"fmt"
	"reflect"
	"sort"
)

var empty struct{}

type IntSet struct {
	ints map[int]struct{}
}

func NewIntSet(ints []int) IntSet {
	s := IntSet{make(map[int]struct{})}
	if ints == nil {
		return s
	}
	for _, n := range ints {
		s.ints[n] = empty
	}
	return s
}

func (s IntSet) Contains(num int) bool {
	if _, ok := s.ints[num]; ok {
		return true
	}
	return false
}

func (s IntSet) Ints() []int {
	var ints []int
	for n := range s.ints {
		ints = append(ints, n)
	}
	sort.Ints(ints)
	return ints
}

func Union(set1, set2 IntSet) IntSet {
	union := NewIntSet(nil)
	for n := range set1.ints {
		union.ints[n] = empty
	}
	for n := range set2.ints {
		union.ints[n] = empty
	}
	return union
}

func Intersection(set1, set2 IntSet) IntSet {
	intersect := NewIntSet(nil)
	for n := range set1.ints {
		if set2.Contains(n) {
			intersect.ints[n] = empty
		}
	}
	return intersect
}

func Equals(set1, set2 IntSet) bool {
	return reflect.DeepEqual(set1, set2)
}

func (s IntSet) String() string {
	return fmt.Sprintf("%v", s.Ints())
}
