package main

import (
	"github.com/at7as/ev7ab/pkg/app"
)

func main() {

	app.Run()

}

// import (
// 	"github.com/at7as/ev7ab/pkg/lab"
// 	"github.com/at7as/ev7ab/pkg/lib"
// )

// func main() {

// 	cfg := lab.Config{
// 		In:     2,
// 		Out:    1,
// 		Target: []float64{0.0},
// 		Limit:  []float64{0.01},
// 		Goal:   true,
// 		Size:   1000,
// 	}
// 	l := lab.New(cfg, &lib.Example1{})
// 	l.Examine()

// 	// l := lab.NewLab(1000, 0.05, 0.2, 0.3, 0.2, 0.5, &producer{}, 1, true)
// 	// l.Add(2, "")
// 	// l.Add(10, "")
// 	// l.Add(1, "")
// 	// l.Examine(20)

// }
