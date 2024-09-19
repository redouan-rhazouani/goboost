package vector

import (
	"slices"
	"sort"
	"testing"
)

// Correctly partitions array with distinct elements
func TestPartitionWithDistinctElements(t *testing.T) {
	a := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	target := 5
	expected := 6 // Assuming the function returns the partition index
	predicate := func(i int) bool { return a[i] < target }
	result := Partition(a, predicate)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
	pp := PartitionPoint(a, predicate)
	if pp != expected {
		t.Errorf("Expected %d, but got %d", expected, pp)
	}
	if !IsPartitioned(len(a), predicate) {
		t.Errorf("Expected %v to be partitioned", a)
	}
}

func TestPartitionAllElementsSatisfyPredicate(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	pred := func(i int) bool { return a[i] > 0 }

	index := Partition(a, pred)

	if index != len(a) {
		t.Errorf("Expected index %d, but got %d", len(a), index)
	}

	pp := PartitionPoint(a, pred)
	if pp != index {
		t.Errorf("Expected %d, but got %d", index, pp)
	}
	if !IsPartitioned(len(a), pred) {
		t.Errorf("Expected %v to be partitioned", a)
	}
}

func TestPartitionAllElementsNotSatisfyingPredicate(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	pred := func(i int) bool { return a[i] < 0 }
	index := Partition(a, pred)

	if index != 0 {
		t.Errorf("Expected index %d, but got %d", 0, index)
	}
	pp := PartitionPoint(a, pred)
	if pp != index {
		t.Errorf("Expected %d, but got %d", index, pp)
	}
	if !IsPartitioned(len(a), pred) {
		t.Errorf("Expected %v to be partitioned", a)
	}
}

func TestPartitionAllIdenticalElements(t *testing.T) {
	a := []int{1, 1, 1, 1, 1}
	pred := func(i int) bool { return a[i] == 1 }
	index := Partition(a, pred)

	if index != len(a) {
		t.Errorf("Expected index %d, but got %d", len(a), index)
	}
	pp := PartitionPoint(a, pred)
	if pp != index {
		t.Errorf("Expected %d, but got %d", index, pp)
	}
	if !IsPartitioned(len(a), pred) {
		t.Errorf("Expected %v to be partitioned", a)
	}
}

func TestPartitionPointSimplePredicate(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	isEven := func(i int) bool { return a[i]%2 == 0 }
	if IsPartitioned(len(a), isEven) {
		t.Errorf("Expected %v to be not partitioned", a)
	}
	Partition(a, isEven)
	result := PartitionPoint(a, isEven)
	expected := 4
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
	if !IsPartitioned(len(a), isEven) {
		t.Errorf("Expected %v to be partitioned", a)
	}
	sort.Sort(sort.IntSlice(a[:result]))
	if firstPart := []int{2, 4, 6, 8}; !slices.Equal(a[:result], firstPart) {
		t.Errorf("First partition Expected %d, but got %d", firstPart, a[:result])
	}
	sort.Sort(sort.IntSlice(a[result:]))
	if secondPart := []int{1, 3, 5, 7, 9}; !slices.Equal(a[result:], secondPart) {
		t.Errorf("Second partition Expected %d, but got %d", secondPart, a[result:])
	}
}
