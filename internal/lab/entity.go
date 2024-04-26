package lab

// type Genome struct { // config for state of lab
// 	design []layout
// 	volume [][][]float64
// }

// type layout struct {
// 	ctrl [3]string
// 	out  [][][2]int
// }

// func (e *entity) compile() {
// e.model = make([]*grid, len(e.design))
// for il, l := range e.design {
// 	g := grid{}
// 	g.aggregator = e.lab.gen.aggr[l.ctrl[0]]
// 	g.processor = e.lab.gen.proc[l.ctrl[1]]
// 	g.consolidator = e.lab.gen.cons[l.ctrl[2]]
// 	g.stack = make([]*atom, len(l.out))
// 	for io, o := range l.out {
// 		a := atom{}
// 		switch l.ctrl[2] {
// 		case "input":
// 			a.in = make([]float64, 1)
// 			a.out = make([]*atom, len(o))
// 		case "link":
// 			a.in = make([]float64, 1)
// 			a.out = make([]*atom, 1)
// 		case "node":
// 			a.out = make([]*atom, len(o))
// 		case "output":
// 			a.in = make([]float64, 1)
// 			a.value = make([]float64, 1)
// 		}
// 		g.stack[io] = &a
// 	}
// 	e.model[il] = &g
// }
// for ig, g := range e.model {
// 	for ia, a := range g.stack {
// 		for i, ref := range e.design[ig].out[ia] {
// 			a.out[i] = e.model[ref[0]].stack[ref[1]]
// 		}
// 	}
// }
// for ig, g := range e.volume {
// 	for ia, a := range g {
// 		for _, v := range a {
// 			e.model[ig].stack[ia].value = append(e.model[ig].stack[ia].value, v)
// 		}
// 	}
// }
// }

// func (e *entity) value(input []float64) []float64 {
// if len(e.model) == 0 {
// 	e.compile()
// }
// ingrid := e.model[0]
// for i, v := range input {
// 	ingrid.stack[i].in[0] = v
// }
// for _, g := range e.model {
// 	for _, a := range g.stack {
// 		g.consolidator.exec(a, g.aggregator, g.processor)
// 	}
// }
// outgrid := e.model[len(e.model)-1]
// output := make([]float64, len(outgrid.stack))
// for i, a := range outgrid.stack {
// 	output[i] = a.value[0]
// }
// 	return nil // output
// }

// func (e *entity) replicate() *entity {
// 	c := entity{
// 		// pt: e.pt,
// 	}
// 	return &c
// }

// func newLayout(a, p, c string, size int) layout {
// 	return layout{
// 		ctrl: [3]string{a, p, c},
// 		out:  make([][][2]int, size),
// 	}
// }
