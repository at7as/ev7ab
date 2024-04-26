package lab

type project struct {
	origin *model
	level  int
	pop    []*entity
	// top    []*entity
	result [][]float64
}

type model struct {
	pt   *project
	nest []*stage
}

type stage struct {
	layout
	Aggregator
	Processor
}

type layout struct {
	a  string
	p  string
	in []int
	n  int
}

type entity struct {
	origin *model
	link   [][][]float64
	node   [][]*atom
	result [][]float64
}

type atom struct {
	in    []float64
	value float64
}

func (p *project) last() []float64 {

	return p.result[len(p.result)-1]
}

func (e *entity) exec(in []float64) *entity {

	for i, s := range e.origin.nest {

		n := e.origin.nest[i].n
		for ii, input := range e.link[i] {

			ini := e.origin.nest[i].in[ii]
			nn := len(input) / e.origin.nest[i].n
			if ini > 0 {

				for iii := 0; iii < nn; iii++ {
					for iiii := 0; iiii < n; iiii++ {
						e.node[i][iiii].in = append(e.node[i][iiii].in, input[iii*n+iiii]*e.node[ini][iii].value)
					}
				}

			} else {

				for iii := 0; iii < nn; iii++ {
					for iiii := 0; iiii < n; iiii++ {
						e.node[i][iiii].in = append(e.node[i][iiii].in, input[iii*n+iiii]*in[iii])
					}
				}

			}

		}

		for _, a := range e.node[i] {
			a.value = s.Processor.exec(a.value, s.Aggregator.exec(a.in))
		}

	}

	return e
}

func (e *entity) value() []float64 {

	out := make([]float64, len(e.node[len(e.node)-1]))
	for i, a := range e.node[len(e.node)-1] {
		out[i] = a.value
	}

	return out
}

func (e *entity) last() []float64 {

	return e.result[len(e.result)-1]
}

// // Entity ...
// type Entity struct {
// 	origin *mind
// 	node   *volume
// 	last   struct {
// 		x []float64
// 		y []float64
// 	}
// 	result []float64
// }

// func newEntity(c *config, v *volume) *Entity {
// 	e := &Entity{origin: c.primal}
// 	if v == nil {
// 		v = newVolume(c.primal, c.rand)
// 	}
// 	e.node = v
// 	return e
// }

// func cloneEntity(v *Entity) *Entity {
// 	e := &Entity{origin: v.origin}
// 	e.node = cloneVolume(v.node)
// 	return e
// }

// // Eval ...
// func (e *Entity) Eval(input []float64) []float64 {
// 	e.last.x = input
// 	output := []float64{}
// 	for i, v := range *e.node {
// 		output = make([]float64, e.origin.output[i].output)
// 		for ii, vv := range v {
// 			r := 0.0
// 			for iii, vvv := range vv {
// 				r += input[iii] * vvv
// 			}
// 			output[ii] = e.origin.output[i].eval.evaluate(r)
// 		}
// 		input = output
// 	}
// 	e.last.y = output
// 	return output
// }

// // Last ...
// func (e *Entity) Last() ([]float64, []float64) {
// 	return e.last.x, e.last.y
// }

// type volume [][][]float64
// func newVolume(m *mind, c *chaos) *volume {
// 	v := make(volume, len(m.output))
// 	for i, a := range m.output {
// 		v[i] = make([][]float64, a.output)
// 		for ii := range v[i] {
// 			v[i][ii] = make([]float64, a.input)
// 			for iii := range v[i][ii] {
// 				v[i][ii][iii] = c.rand101()
// 			}
// 		}
// 	}
// 	return &v
// }

// func cloneVolume(value *volume) *volume {
// 	v := make(volume, len(*value))
// 	for i, a := range *value {
// 		v[i] = make([][]float64, len(a))
// 		for ii, aa := range a {
// 			v[i][ii] = make([]float64, len(aa))
// 			copy(v[i][ii], (*value)[i][ii])
// 		}
// 	}
// 	return &v
// }

// type area struct {
// 	input  int
// 	output int
// 	eval   evaluater
// }

// type mind struct {
// 	input  int
// 	output []*area
// }

// func (m *mind) add(node int, evaluater string) {
// 	if m.input == 0 {
// 		m.input = node
// 	} else {
// 		input := 0
// 		if len(m.output) == 0 {
// 			input = m.input
// 		} else {
// 			input = m.output[len(m.output)-1].output
// 		}
// 		m.output = append(m.output, &area{input, node, getEvaluater(evaluater)})
// 	}
// }
