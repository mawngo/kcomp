package kmeans

import (
	"errors"
	"math"
)

var (
	ErrEmptySet       = errors.New("empty training set")
	ErrZeroIterations = errors.New("number of iterations cannot be less than 1")
	ErrOneCluster     = errors.New("number of clusters cannot be less than 2")
)

// DistanceFunc represents a function for measuring distance
// between n-dimensional vectors.
type DistanceFunc func([]float64, []float64) float64

// Online represents parameters important for online learning in
// clustering algorithms.
type Online struct {
	Alpha     float64
	Dimension int
}

// HCEvent represents the intermediate result of computation of hard clustering algorithm
// and are transmitted periodically to the caller during online learning
type HCEvent struct {
	Cluster     int
	Observation []float64
}

var (
	// EuclideanDistance is one of the common distance measurement
	EuclideanDistance = func(a, b []float64) float64 {
		var (
			s, t float64
		)

		for i, _ := range a {
			t = a[i] - b[i]
			s += t * t
		}

		return math.Sqrt(s)
	}

	// EuclideanDistanceSquared is one of the common distance measurement
	EuclideanDistanceSquared = func(a, b []float64) float64 {
		var (
			s, t float64
		)

		for i, _ := range a {
			t = a[i] - b[i]
			s += t * t
		}

		return s
	}
)
