package intset

import (
	"reflect"
	"testing"
)

func TestNewIntSet(t *testing.T) {
	expected := IntSet{
		ints: map[int]struct{}{
			1:  empty,
			2:  empty,
			5:  empty,
			-3: empty,
		},
	}
	result := NewIntSet([]int{1, 2, 1, 5, -3})
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expecting %s, got %s", expected, result)
	}
}

func TestContains(t *testing.T) {
	ints := []int{1, 2, 3}
	s := NewIntSet(ints)
	if !s.Contains(1) {
		t.Errorf("set should contain 1, got %s", s)
	}
}

func TestUnion(t *testing.T) {
	set1 := NewIntSet([]int{1, 4, 7})
	set2 := NewIntSet([]int{-4, 7, 10})
	expected := NewIntSet([]int{-4, 1, 4, 7, 10})
	result := Union(set1, set2)
	if !Equals(result, expected) {
		t.Errorf("expecting %s, got %s", expected, result)
	}
}

func TestIntersection(t *testing.T) {
	set1 := NewIntSet([]int{1, 4, 7})
	set2 := NewIntSet([]int{-4, 7, 10})
	result := Intersection(set1, set2)
	if len(result.ints) > 1 {
		t.Errorf("expecting 1 integer, found %d", len(result.ints))
	}
	if !result.Contains(7) {
		t.Errorf("expecting the set to contain the single int 7, got %s", result)
	}
}
