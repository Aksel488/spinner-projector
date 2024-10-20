package util

import (
	"fmt"
	"image"
	"math"
	"strconv"
	"strings"
)

func Distance(x1, y1, x2, y2 int) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(float64((a * a) + (b * b)))
}

func PDistance(start, end image.Point) float64 {
	return Distance(start.X, start.Y, end.X, end.Y)
}

// CountDecimals returns the number of decimal places in the given float32
func CountDecimals(f float32) int {
	str := strconv.FormatFloat(float64(f), 'f', -1, 32)
	parts := strings.Split(str, ".")
	if len(parts) == 2 {
		return len(parts[1])
	}
	return 0
}

// GetLeftPadding calculates the number of characters based on the length of the max value
func GetLeftPadding(max float32) int {
	extra := 0
	if float64(max) > math.Abs(float64(max)) {
		extra = 1
	}
	str := strconv.FormatFloat(float64(max), 'f', -1, 32)
	parts := strings.Split(str, ".")
	return len(parts[0]) + extra
}

// FormatValue formats the value with dynamic decimals and left padding
func FormatLogValue(value float32, minDecimals, leftPadding int) string {
	format := "%0" + strconv.Itoa(leftPadding+minDecimals+1) + "." + strconv.Itoa(minDecimals) + "f"
	return fmt.Sprintf(format, value)
}
