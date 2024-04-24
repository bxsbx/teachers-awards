package util

import "math"

// 保留n位小数
func Decimals(x float64, n int) float64 {
	pow10 := math.Pow10(n)
	return math.Round(x*pow10) / pow10
}

// 同比
func OnYear(a, b int64) float64 {
	if b > 0 {
		return Decimals(float64(a-b)/float64(b)*100, 2)
	} else if a > 0 {
		return 100
	}
	return 0
}

// 占比
func Rate(a, b int64) float64 {
	if b > 0 {
		return Decimals(float64(a)/float64(b)*100, 2)
	}
	return 0
}
