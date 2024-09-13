package app

import (
	"slices"

	"github.com/at7as/ev7ab/pkg/lab"
)

type model struct {
	model  []*model
	size   int
	volume int
	source []position
	valid  bool
}

func newModel(size int) *model {

	return &model{model: make([]*model, 0, size)}
}

func newNode(size int, source []position) *model {

	if source == nil {
		source = make([]position, 0)
	}

	return &model{size: size, source: source}
}

func newProjectModel(o *project, in []int, out int) *model {

	m := newModel(2)
	if o != nil {
		m = o.model.clone()
	} else {
		m.addModel(nil)
		m.addModel(nil)
		for _, i := range in {
			m.model[0].addModel(newNode(i, nil))
		}
		m.model[1].addModel(newNode(out, []position{{0, 0}}))
	}
	m.measure()

	return m
}

func createProjectModel(layout [][]lab.Node) *model {

	m := newModel(len(layout))
	for i := range len(layout) {
		s := newModel(len(layout[i]))
		for ii := range len(layout[i]) {
			n := newNode(layout[i][ii].Out, nil)
			for _, pos := range layout[i][ii].Src {
				n.source = append(n.source, newPosition(pos[0], pos[1]))
			}
			s.addModel(n)
		}
		m.addModel(s)
	}

	return m
}

func (m *model) clone() *model {

	d := newModel(len(m.model))
	for _, s := range m.model {
		ds := newModel(len(s.model))
		for _, n := range s.model {
			dn := newNode(n.size, nil)
			for _, src := range n.source {
				dn.source = append(dn.source, src)
			}
			ds.addModel(dn)
		}
		d.addModel(ds)
	}

	return d
}

func (m *model) addModel(n *model) {

	if n == nil {
		n = newModel(0)
	}

	m.model = append(m.model, n)

}

func (m *model) insertModel(n *model, index int) {

	if n == nil {
		n = newModel(0)
	}

	m.model = slices.Insert(m.model, index, n)

}

func (m *model) deleteModel(index int) {

	m.model = slices.Delete(m.model, index, index+1)

}

func (m *model) measure() {

	m.size = 0
	m.volume = 0
	for i := range m.model {
		m.model[i].size = 0
		m.model[i].volume = 0
		for ii := range m.model[i].model {
			m.model[i].model[ii].volume = 0
		}
	}

	for i, s := range m.model {
		for ii, n := range s.model {
			m.model[i].size += n.size
			m.size += n.size
			for _, src := range n.source {
				v := m.model[src.x].model[src.y].size * n.size
				m.model[i].model[ii].volume += v
				m.model[i].volume += v
				m.volume += v
			}
		}
	}

}

func (m *model) validate() bool {

	for i := range m.model {
		for ii := range m.model[i].model {
			m.model[i].model[ii].valid = i == len(m.model)-1
		}
	}

	for i := range m.model {
		for ii := range m.model[i].model {
			for _, src := range m.model[i].model[ii].source {
				m.model[src.x].model[src.y].valid = true
			}
		}
	}

	for i := range m.model {
		for ii := range m.model[i].model {
			if !m.model[i].model[ii].valid {
				continue
			}
			m.model[i].model[ii].valid = m.model[i].model[ii].size > 0
			if !m.model[i].model[ii].valid {
				continue
			}
			if i > 0 {
				m.model[i].model[ii].valid = len(m.model[i].model[ii].source) > 0
			}
		}
	}

	for _, s := range m.model {
		if len(s.model) == 0 {
			return false
		}
		for _, n := range s.model {
			if !n.valid {
				return false
			}
		}
	}

	return true
}

func (m *model) insertStage(index int) {

	m.insertModel(nil, index)
	for _, s := range m.model {
		for _, n := range s.model {
			for iii, src := range n.source {
				if index <= src.x {
					n.source[iii].x++
				}
			}
		}
	}

}

func (m *model) deleteStage(index int) {

	if index != 0 && index != len(m.model)-1 {
		m.deleteModel(index)
		for _, s := range m.model {
			for _, n := range s.model {
				for iii, src := range n.source {
					if index <= src.x {
						n.source[iii].x--
					}
				}
			}
		}
	}

}

func (m *model) convert() [][]lab.Node {

	layout := make([][]lab.Node, len(m.model))

	for i, s := range m.model {
		layout[i] = make([]lab.Node, len(s.model))
		for ii, n := range s.model {
			src := make([][2]int, len(n.source))
			for iii, pos := range n.source {
				src[iii] = [2]int{pos.x, pos.y}
			}
			layout[i][ii] = lab.Node{Src: src, Out: n.size}
		}
	}

	return layout
}
