package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunked(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, [][]int{{1}, {2}, {3}, {4}, {5}, {6}}, SplitToBulksInt(s, 1))
	assert.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, SplitToBulksInt(s, 100))
	assert.Equal(t, [][]int{{1, 2, 3, 4, 5, 6}}, SplitToBulksInt(s, 6))
	assert.Equal(t, [][]int{{1, 2, 3, 4, 5}, {6}}, SplitToBulksInt(s, 5))
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, SplitToBulksInt(s, 3))
	assert.Equal(t, [][]int{}, SplitToBulksInt(s, 0))
	assert.Equal(t, [][]int{}, SplitToBulksInt([]int{}, 4))
}

func TestReverseMap(t *testing.T) {
	assert.Equal(
		t,
		map[int]int{11: 1, 22: 2},
		ReverseMapIntToInt(map[int]int{1: 11, 2: 22}),
	)

}

func TestRemoveElements(t *testing.T) {
	s := []int{10, 20, 30, 40, 50, 60}
	assert.Equal(t, []int{}, RemoveElementsInt(s, s))
	assert.Equal(t, []int{20, 30, 40, 50, 60}, RemoveElementsInt(s, []int{10}))
	assert.Equal(t, []int{40, 50, 60}, RemoveElementsInt(s, []int{10, 20, 30, 100, 200}))
	assert.Equal(t, s, RemoveElementsInt(s, []int{100, 200}))
}
