package kmeans

import (
	"math"
)

// DistanceFunc represents a function for measuring distance between n-dimensional vectors.
type DistanceFunc func([]float64, []float64) float64

var (
	// EuclideanDistance is one of the common distance measurement.
	EuclideanDistance = func(a, b []float64) float64 {
		var (
			s, t float64
		)

		for i := range a {
			t = a[i] - b[i]
			s += t * t
		}

		return math.Sqrt(s)
	}

	// EuclideanDistanceSquared is one of the common distance measurement.
	EuclideanDistanceSquared = func(a, b []float64) float64 {
		var (
			s, t float64
		)

		for i := range a {
			t = a[i] - b[i]
			s += t * t
		}

		return s
	}
)
