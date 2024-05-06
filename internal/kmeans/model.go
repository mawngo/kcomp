package kmeans

import (
	"gonum.org/v1/gonum/floats"
	"math"
	"math/rand"
	"runtime"
)

type Dataset [][]float64

type Trainer struct {
	k             int
	maxIterations int
	distanceFn    DistanceFunc
	delta         float64
	concurrency   int
}

type TrainerOption func(*Trainer)

type Model struct {
	distanceFn DistanceFunc
	k          int
	data       Dataset
	centroids  Dataset
	mapping    []int
	iter       int
}

// NewTrainer create new Trainer
func NewTrainer(k int, options ...TrainerOption) Trainer {
	t := Trainer{
		k:             k,
		maxIterations: 100,
		distanceFn:    EuclideanDistance,
		delta:         0.01,
		concurrency:   runtime.NumCPU(),
	}
	for i := range options {
		options[i](&t)
	}
	return t
}

func WithDistanceFunc(fn DistanceFunc) TrainerOption {
	return func(t *Trainer) {
		t.distanceFn = fn
	}
}

func WithMaxIterations(i int) TrainerOption {
	return func(t *Trainer) {
		t.maxIterations = i
	}
}

func WithDeltaThreshold(delta float64) TrainerOption {
	return func(t *Trainer) {
		t.delta = delta
	}
}

// Fit create and train the *Model.
func (t Trainer) Fit(data Dataset) *Model {
	model := Model{data: data, k: t.k, distanceFn: t.distanceFn}
	model.initializeMean()
	l := len(model.centroids[0])
	changeThreshold := int(float64(len(data)) * t.delta)

	cb, cn := prepare(t.k, l)
	iter := 0
	for ; iter < t.maxIterations; iter++ {
		changes := 0
		icb := make([][]int, t.concurrency)
		icn := make([]Dataset, t.concurrency)
		ch := make(chan int, t.concurrency)
		for num := range t.concurrency {
			go func() {
				defer func() {
					ch <- num
				}()
				cb, cn := prepare(t.k, l)
				for i := num; i < len(data); i += t.concurrency {
					m := t.distanceFn(data[i], model.centroids[0])
					n := 0

					for j := 1; j < t.k; j++ {
						if d := t.distanceFn(data[i], model.centroids[j]); d < m {
							m = d
							n = j
						}
					}

					if model.mapping[i] != n {
						changes++
					}

					model.mapping[i] = n
					cb[n]++
					floats.Add(cn[n], data[i])
				}
				icb[num] = cb
				icn[num] = cn
			}()
		}

		for range t.concurrency {
			num := <-ch
			for n := range t.k {
				cb[n] += icb[num][n]
				floats.Add(cn[n], icn[num][n])
			}
		}

		for i := 0; i < t.k; i++ {
			floats.Scale(1/float64(cb[i]), cn[i])
			cb[i] = 0

			for j := 0; j < l; j++ {
				model.centroids[i][j] = cn[i][j]
				cn[i][j] = 0
			}
		}

		if changes < changeThreshold {
			break
		}
	}

	model.iter = iter
	return &model
}

func prepare(k int, l int) ([]int, Dataset) {
	cb := make([]int, k)
	cn := make(Dataset, k)
	for i := 0; i < k; i++ {
		cn[i] = make([]float64, l)
	}
	return cb, cn
}

func (m *Model) initializeMean() {
	m.mapping = make([]int, len(m.data))
	m.centroids = make(Dataset, m.k)
	m.centroids[0] = m.data[rand.Intn(len(m.data)-1)]

	d := make([]float64, len(m.data))
	for i := 1; i < m.k; i++ {
		s := float64(0)
		for j := 0; j < len(m.data); j++ {
			l := m.distanceFn(m.centroids[0], m.data[j])
			for g := 1; g < i; g++ {
				if f := m.distanceFn(m.centroids[g], m.data[j]); f < l {
					l = f
				}
			}

			d[j] = math.Pow(l, 2)
			s += d[j]
		}

		t := rand.Float64() * s
		k := 0
		for s = d[0]; s < t; s += d[k] {
			k++
		}

		m.centroids[i] = m.data[k]
	}
}

// Predict returns number of cluster to which the observation would be assigned.
func (m *Model) Predict(p []float64) int {
	l := 0
	n := m.distanceFn(p, m.centroids[0])
	for i := 1; i < m.k; i++ {
		if d := m.distanceFn(p, m.centroids[i]); d < n {
			n = d
			l = i
		}
	}
	return l
}

// Guesses returns mapping from data point indices to cluster numbers.
func (m *Model) Guesses() []int {
	return m.mapping
}

// Cluster returns cluster at position i.
func (m *Model) Cluster(i int) []float64 {
	return m.centroids[i]
}

// Iter returns model number of iterations.
func (m *Model) Iter() int {
	return m.iter
}
