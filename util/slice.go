package util

import "strconv"

// InSliceInt checks if needle is in the stack
func InSliceInt(needle int, stack []int) bool {
	for _, item := range stack {
		if needle == item {
			return true
		}
	}
	return false
}

// SliceIntDiff differences s1 against s2
func SliceIntDiff(s1 []int, s2 []int) []int {
	var result []int
	for _, val := range s1 {
		if !InSliceInt(val, s2) {
			result = append(result, val)
		}
	}
	return result
}

// SliceIntToString convert []int to []string
func SliceIntToString(s []int) []string {
	var r = make([]string, 0, len(s))
	for _, item := range s {
		r = append(r, strconv.Itoa(item))
	}
	return r
}
