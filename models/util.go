package models

import (
	"image"
	"math"
)

func Distance(x1, y1, x2, y2 int) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(float64((a * a) + (b * b)))
}

func PDistance(start, end image.Point) float64 {
	return Distance(start.X, start.Y, end.X, end.Y)
}
