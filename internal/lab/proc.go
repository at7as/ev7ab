package lab

import "math"

// Processor ...
type Processor func(float64, float64) float64

// (-*,*)
func linearProc(v, in float64) float64 {

	return in
}

// {0,1}
func binaryProc(v, in float64) float64 {

	if in < 0 {
		return 0.0
	}

	return 1.0
}

// [0,*)
func reluProc(v, in float64) float64 {

	if in < 0 {
		return 0.0
	}

	return in
}

// (-*,*)
func lreluProc(v, in float64) float64 {

	if in < 0 {
		return in * 0.01
	}

	return in
}

// [-0.278,*)
func siluProc(v, in float64) float64 {

	return in / (1.0 + math.Exp(-in))
}

// (0,1)
func sigmoidProc(v, in float64) float64 {

	return 1.0 / (1.0 + math.Exp(-in))
}

// (0,*)
func softplusProc(v, in float64) float64 {

	return math.Log(1.0 + math.Exp(in))
}

// (-1,1)
func softsignProc(v, in float64) float64 {

	return in / (1.0 + math.Abs(in))
}

// (-1,1)
func tanhProc(v, in float64) float64 {

	return math.Tanh(in)
}

// (-pi/2,pi/2)
func arctanProc(v, in float64) float64 {

	return math.Atan(in)
}

// [-1,1]
func sinusoidProc(v, in float64) float64 {

	return math.Sin(in)
}

// (0,1]
func gaussianProc(v, in float64) float64 {

	return math.Exp(-math.Pow(in, 2.0))
}
