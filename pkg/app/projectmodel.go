package app

type projectModel struct {
	model []*projectModelStage
	// dirty  bool
	size   int
	volume int
}

type projectModelStage struct {
	stage  []*projectModelNode
	size   int
	volume int
}

type projectModelNode struct {
	source []projectModelSource
	size   int
	volume int
	valid  bool
}

type projectModelSource struct {
	s, n int
}

func newProjectModel(o *project, in, out int) *projectModel {

	m := &projectModel{}
	if o != nil {
		m = cloneModel(o.m)
	} else {
		m.model = make([]*projectModelStage, 2)
		m.model[0] = &projectModelStage{stage: make([]*projectModelNode, 1)}
		m.model[0].stage[0] = &projectModelNode{source: make([]projectModelSource, 0), size: in}
		m.model[1] = &projectModelStage{stage: make([]*projectModelNode, 1)}
		m.model[1].stage[0] = &projectModelNode{source: []projectModelSource{{0, 0}}, size: out}
	}
	// m.dirty = true
	m.measure()

	return m
}

func (m *projectModel) measure() {

	m.size = 0
	m.volume = 0
	for i := range m.model {
		m.model[i].size = 0
		m.model[i].volume = 0
		for ii := range m.model[i].stage {
			m.model[i].stage[ii].volume = 0
		}
	}

	for i, s := range m.model {
		for ii, n := range s.stage {
			m.model[i].size += n.size
			m.size += n.size
			for _, src := range n.source {
				v := m.model[src.s].stage[src.n].size * n.size
				m.model[i].stage[ii].volume += v
				m.model[i].volume += v
				m.volume += v
			}
		}
	}

}

func (m *projectModel) addStage(index int) {

	m.model = append(m.model[:index+1], m.model[index:]...)
	m.model[index] = &projectModelStage{stage: make([]*projectModelNode, 0)}
	for i := range m.model {
		for ii := range m.model[i].stage {
			for iii, src := range m.model[i].stage[ii].source {
				if index <= src.s {
					m.model[i].stage[ii].source[iii].s++
				}
			}
		}
	}

}

func (m *projectModel) removeStage(index int) {

	if index != 0 && index != len(m.model)-1 {
		m.model = append(m.model[:index], m.model[index+1:]...)
		for i := range m.model {
			for ii := range m.model[i].stage {
				for iii, src := range m.model[i].stage[ii].source {
					if index <= src.s {
						m.model[i].stage[ii].source[iii].s--
					}
				}
			}
		}
	}

}

func (s *projectModelStage) addNode() {

	s.stage = append(s.stage, &projectModelNode{source: make([]projectModelSource, 0)})
}

func (s *projectModelStage) removeNode(index int) {

	if len(s.stage) > index {
		s.stage = append(s.stage[:index], s.stage[index+1:]...)
	}

}

func cloneModel(m *projectModel) *projectModel {

	n := &projectModel{}
	n.model = make([]*projectModelStage, len(m.model))
	for i, s := range m.model {
		n.model[i] = &projectModelStage{stage: make([]*projectModelNode, len(s.stage))}
		for ii, node := range s.stage {
			n.model[i].stage[ii] = &projectModelNode{
				source: make([]projectModelSource, len(node.source)),
				size:   node.size,
			}
			for iii, src := range node.source {
				n.model[i].stage[ii].source[iii] = projectModelSource{src.s, src.n}
			}
		}
	}

	return n
}
