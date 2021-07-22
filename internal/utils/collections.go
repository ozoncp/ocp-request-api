package utils

import (
	"math"
)

// SplitToBulksInt Converts a given slice of ints to a slice of slices of ints of a given size.
func SplitToBulksInt(items []int, chunkSize int) [][]int {
	if chunkSize <= 0 {
		return make([][]int, 0)
	}

	chunksNum := int(math.Ceil(float64(len(items)) / float64(chunkSize)))

	ret := make([][]int, 0, chunksNum)

	for chunkStart := 0; chunkStart < len(items); chunkStart = chunkStart + chunkSize {
		chunkEnd := chunkStart + chunkSize
		if chunkEnd > len(items) {
			chunkEnd = len(items)
		}

		ret = append(ret, items[chunkStart:chunkEnd])
	}
	return ret
}

// ReverseMapIntToInt Converts mapping to a reversed mapping (a map where key becomes a value and vice-versa).
func ReverseMapIntToInt(mapping map[int]int) map[int]int {
	reversed := make(map[int]int, len(mapping))

	for key, val := range mapping {
		reversed[val] = key
	}
	return reversed
}

// RemoveElementsInt Return a new slice with elements in the `items` that are not in the `remove` slice.
func RemoveElementsInt(items []int, remove []int) []int {
	removeSet := make(map[int]bool, len(remove))
	ret := make([]int, 0, len(items))

	for _, val := range remove {
		removeSet[val] = true
	}

	for _, val := range items {
		if removeSet[val] {
			continue
		}
		ret = append(ret, val)
	}

	return ret
}
