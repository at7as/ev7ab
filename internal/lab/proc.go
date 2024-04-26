package lab

import "math"

type Processor interface {
	exec(float64, float64) float64
}

// (-*,*)
type linearProc struct{}

func (proc *linearProc) exec(v, in float64) float64 {

	return in
}

// {0,1}
type binaryProc struct{}

func (proc *binaryProc) exec(v, in float64) float64 {

	if in < 0 {
		return 0.0
	}

	return 1.0
}

// [0,*)
type reLUProc struct{}

func (proc *reLUProc) exec(v, in float64) float64 {

	if in < 0 {
		return 0.0
	}

	return in
}

// (-*,*)
type lReLUProc struct{}

func (proc *lReLUProc) exec(v, in float64) float64 {

	if in < 0 {
		return in * 0.01
	}

	return in
}

// [-0.278,*)
type siLUProc struct{}

func (proc *siLUProc) exec(v, in float64) float64 {

	return in / (1.0 + math.Exp(-in))
}

// (0,1)
type sigmoidProc struct{}

func (proc *sigmoidProc) exec(v, in float64) float64 {

	return 1.0 / (1.0 + math.Exp(-in))
}

// (0,*)
type softPlusProc struct{}

func (proc *softPlusProc) exec(v, in float64) float64 {

	return math.Log(1.0 + math.Exp(in))
}

// (-1,1)
type softSignProc struct{}

func (proc *softSignProc) exec(v, in float64) float64 {

	return in / (1.0 + math.Abs(in))
}

// (-1,1)
type tanhProc struct{}

func (proc *tanhProc) exec(v, in float64) float64 {

	return math.Tanh(in)
}

// (-pi/2,pi/2)
type arcTanProc struct{}

func (proc *arcTanProc) exec(v, in float64) float64 {

	return math.Atan(in)
}

// [-1,1]
type sinusoidProc struct{}

func (proc *sinusoidProc) exec(v, in float64) float64 {

	return math.Sin(in)
}

// (0,1]
type gaussianProc struct{}

func (proc *gaussianProc) exec(v, in float64) float64 {

	return math.Exp(-math.Pow(in, 2.0))
}
