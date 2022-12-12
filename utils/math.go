package utils

import "math"

const INT_MAX = int(^uint(0) >> 1)

const INT_MIN = ^INT_MAX

func GetFloat(val interface{}) (float64, bool) {
	if val1, ok := val.(float64); ok {
		return val1, true
	}
	if val2, ok := val.(int32); ok {
		return float64(val2), true
	}
	if val3, ok := val.(int64); ok {
		return float64(val3), true
	}
	return 0, false
}

func Zero(nums ...float64) bool {
	for _, v := range nums {
		if math.Abs(v) < 1e-6 {
			return true
		}
	}
	return false
}

func ZeroList(l []float64) bool {
	for _, v := range l {
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

func Bool2Int(b bool) int {
	if b {
		return 1
	} else {
		return 0
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
