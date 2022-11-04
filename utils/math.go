package utils

import "math"

const INT_MAX = int(^uint(0) >> 1)

const INT_MIN = ^INT_MAX

func Zero(nums ...float64) bool {
	for _, v := range nums {
		if math.Abs(v) < 1e-6 {
			return true
		}
	}
	return false
}

func Bool2Float(b bool) float64 {
	if b {
		return 1.0
	} else {
		return 0.0
	}
}

func Max(nums ...int) int {
	var ans int = INT_MIN
	for _, v := range nums {
		if v > ans {
			ans = v
		}
	}
	return ans
}

func Min(nums ...int) int {
	var ans int = INT_MAX
	for _, v := range nums {
		if v < ans {
			ans = v
		}
	}
	return ans
}
