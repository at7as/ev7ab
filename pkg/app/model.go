package app

import (
	"slices"
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

func newProjectModel(o *project, in, out int) *model {

	m := newModel(2)
	if o != nil {
		m = o.model.clone()
	} else {
		m.addModel(nil)
		m.addModel(nil)
		m.model[0].addModel(newNode(in, nil))
		m.model[1].addModel(newNode(out, []position{{0, 0}}))
	}
	m.measure()

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
