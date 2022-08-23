package react

var _ Reactor = &Sheet{}

func New() *Sheet {
	s := &Sheet{
		ch:   make(chan Update),
		deps: make(map[Cell][]dependencies),
	}

	return s
}

type Sheet struct {
	Inputs       []*Input
	ComputeCells []*Compute

	deps map[Cell][]dependencies

	ch chan Update
	id int
}

type Update struct {
	id    int
	value int
}

type dependencies struct {
	f1     func(int) int
	f2     func(int, int) int
	target *Compute
	f1Deps []Cell
	f2Deps []Cell
}

func (s *Sheet) CreateInput(i int) InputCell {
	cell := Input{value: i, id: s.id, ch: s.ch, sheet: s}
	s.Inputs = append(s.Inputs, &cell)
	s.id++
	return &cell
}

func (s *Sheet) CreateCompute1(cell Cell, f func(int) int) ComputeCell {
	compute := &Compute{
		sheet: s,
		value: f(cell.Value()),
		f1:    f,
		cells: []Cell{cell},
	}

	s.deps[cell] = append(s.deps[cell], dependencies{
		f1:     f,
		target: compute,
		f1Deps: []Cell{cell},
	})

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}

func (s *Sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	compute := &Compute{
		sheet: s,
		value: f(cell.Value(), cell2.Value()),
		f2:    f,
		cells: []Cell{cell, cell2},
	}

	s.deps[cell] = append(s.deps[cell], dependencies{f2: f, target: compute, f2Deps: []Cell{cell, cell2}})
	s.deps[cell2] = append(s.deps[cell2], dependencies{f2: f, target: compute, f2Deps: []Cell{cell, cell2}})

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}
