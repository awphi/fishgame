package util

import (
	"fmt"
	"math"
	"math/rand"
)

// TODO could be generic? No need for now
type Distribution interface {
	// evaluate the probability density function (PDF) of this distribution for a given value
	Sample(val float64) float64
	// generate a random number in range [-math.MaxFloat64, math.MaxFloat64] weighted by this distribution
	Random() float64
}

type NormalDistribution struct {
	Sd   float64
	Mean float64
}

type UniformDistribution struct {
	Min float64
	Max float64
}

type ExponentialDistribution struct {
	Lambda float64
}

type LogNormalDistribution struct {
	Mean  float64
	Sigma float64
}

type TriangularDistribution struct {
	Min  float64
	Max  float64
	Mode float64
}

func interp(x float64, x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	if x < x1 {
		return y1
	}

	if x > x2 {
		return y2
	}

	return y1 + ((x-x1)/(x2-x1))*(y2-y1)
}

func (d *NormalDistribution) Sample(val float64) float64 {
	t1 := (1 / (d.Sd * math.Sqrt(math.Pi*2)))
	exp := -0.5 * math.Pow(((val-d.Mean)/d.Sd), 2)
	return t1 * math.Pow(math.E, exp)
}

func (d *NormalDistribution) Random() float64 {
	return rand.NormFloat64()*d.Sd + d.Mean
}

func (d *UniformDistribution) Sample(val float64) float64 {
	if val > d.Max || val < d.Min {
		return 0
	}

	r := d.Max - d.Min
	if r == 0 {
		return 1
	}
	return 1 / r
}

func (d *UniformDistribution) Random() float64 {
	r := d.Max - d.Min
	return d.Min + r*rand.Float64()
}

func (d *ExponentialDistribution) Sample(val float64) float64 {
	if val < 0 {
		return 0
	}

	return d.Lambda * math.Pow(math.E, -d.Lambda*val)
}

func (d *ExponentialDistribution) Random() float64 {
	// https://en.wikipedia.org/wiki/Inverse_transform_sampling
	return -math.Log(1-rand.Float64()) / d.Lambda
}

func (d *LogNormalDistribution) Sample(val float64) float64 {
	scalar := 1 / (val * d.Sigma * math.Sqrt(2*math.Pi))
	exp := -(math.Pow(math.Log(val)-d.Mean, 2) / (2 * math.Pow(d.Sigma, 2)))
	return scalar * math.Pow(math.E, exp)
}

func (d *LogNormalDistribution) Random() float64 {
	return math.Pow(math.E, rand.NormFloat64()*d.Sigma+d.Mean)
}

func (d *TriangularDistribution) Sample(val float64) float64 {
	peak := 2 / (d.Max - d.Min)

	if val < d.Mode {
		return interp(val, d.Min, d.Mode, 0, peak)
	}

	return interp(val, d.Mode, d.Max, peak, 0)
}

func (d *TriangularDistribution) Random() float64 {
	u := rand.Float64()
	fc := (d.Max - d.Min) / (d.Mode - d.Min)

	if u < fc {
		return d.Min + math.Sqrt(u*(d.Mode-d.Min)*(d.Max-d.Min))
	} else {
		return d.Mode - math.Sqrt((1-u)*(d.Mode-d.Min)*(d.Mode-d.Max))
	}
}

func NewDistribution(distType string) (Distribution, error) {
	if distType == "normal" {
		return &NormalDistribution{}, nil
	} else if distType == "exponential" {
		return &ExponentialDistribution{}, nil
	} else if distType == "lognormal" {
		return &LogNormalDistribution{}, nil
	} else if distType == "triangular" {
		return &TriangularDistribution{}, nil
	} else if distType == "uniform" {
		return &UniformDistribution{}, nil
	}

	return &NormalDistribution{}, fmt.Errorf("unrecognised distribution type '%s'", distType)
}
