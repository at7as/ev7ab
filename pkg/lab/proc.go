package lab

import "math"

// Processor ...
type Processor func(float64) float64

// (-*,*)
func linearProc(v float64) float64 {

	return v
}

// {0,1}
func binaryProc(v float64) float64 {

	if v < 0 {
		return 0.0
	}

	return 1.0
}

// [0,*)
func reluProc(v float64) float64 {

	if v < 0 {
		return 0.0
	}

	return v
}

// (-*,*)
func lreluProc(v float64) float64 {

	if v < 0 {
		return v * 0.01
	}

	return v
}

// [-0.278,*)
func siluProc(v float64) float64 {

	return v / (1.0 + math.Exp(-v))
}

// (0,1)
func sigmoidProc(v float64) float64 {

	return 1.0 / (1.0 + math.Exp(-v))
}

// (0,*)
func softplusProc(v float64) float64 {

	return math.Log(1.0 + math.Exp(v))
}

// (-1,1)
func softsignProc(v float64) float64 {

	return v / (1.0 + math.Abs(v))
}

// (-1,1)
func tanhProc(v float64) float64 {

	return math.Tanh(v)
}

// (-pi/2,pi/2)
func arctanProc(v float64) float64 {

	return math.Atan(v)
}

// [-1,1]
func sinusoidProc(v float64) float64 {

	return math.Sin(v)
}

// (0,1]
func gaussianProc(v float64) float64 {

	return math.Exp(-math.Pow(v, 2.0))
}

var procMap map[string]Processor = map[string]Processor{
	"linear":   linearProc,
	"binary":   binaryProc,
	"relu":     reluProc,
	"lrelu":    lreluProc,
	"silu":     siluProc,
	"sigmoid":  sigmoidProc,
	"softplus": softplusProc,
	"softsign": softsignProc,
	"tanh":     tanhProc,
	"arctan":   arctanProc,
	"sinusoid": sinusoidProc,
	"gaussian": gaussianProc,
}
