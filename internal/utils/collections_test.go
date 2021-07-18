package utils

import (
	"reflect"
	"testing"
)

func TestChunked(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	assertDeapEqual(t, ChunkedInt(s, 1), [][]int{{1}, {2}, {3}, {4}, {5}, {6}})
	assertDeapEqual(t, ChunkedInt(s, 100), [][]int{{1, 2, 3, 4, 5, 6}})
	assertDeapEqual(t, ChunkedInt(s, 6), [][]int{{1, 2, 3, 4, 5, 6}})
	assertDeapEqual(t, ChunkedInt(s, 5), [][]int{{1, 2, 3, 4, 5}, {6}})
	assertDeapEqual(t, ChunkedInt(s, 3), [][]int{{1, 2, 3}, {4, 5, 6}})
	assertDeapEqual(t, ChunkedInt(s, 0), [][]int{})
	assertDeapEqual(t, ChunkedInt([]int{}, 4), [][]int{})
}

func TestReverseMap(t *testing.T) {
	assertDeapEqual(
		t,
		ReverseMapIntToInt(map[int]int{1: 11, 2: 22}),
		map[int]int{11: 1, 22: 2},
	)

}

func TestRemoveElements(t *testing.T) {
	s := []int{10, 20, 30, 40, 50, 60}
	assertDeapEqual(t, RemoveElementsInt(s, s), []int{})
	assertDeapEqual(t, RemoveElementsInt(s, []int{10}), []int{20, 30, 40, 50, 60})
	assertDeapEqual(t, RemoveElementsInt(s, []int{10, 20, 30, 100, 200}), []int{40, 50, 60})
	assertDeapEqual(t, RemoveElementsInt(s, []int{100, 200}), s)
}

func assertDeapEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected %v. Actual %v", expected, actual)
	}
}
