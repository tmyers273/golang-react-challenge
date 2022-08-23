package react

var _ Reactor = &Sheet{}

func New() *Sheet {
	s := &Sheet{
		ch:   make(chan Update),
		deps: make(map[Cell]dependencies),
	}

	return s
}

type Sheet struct {
	Inputs       []*Input
	ComputeCells []*Compute

	deps map[Cell]dependencies

	ch chan Update
	id int
}

type Update struct {
	id    int
	value int
}

type dependencies struct {
	f1      func(int) int
	f2      func(int, int) int
	targets []Cell
	deps    []Cell
	f2Deps  [][]Cell
}

func (s *Sheet) CreateInput(i int) InputCell {
	cell := Input{value: i, id: s.id, ch: s.ch, sheet: s}
	s.Inputs = append(s.Inputs, &cell)
	s.id++
	return &cell
}

func (s *Sheet) CreateCompute1(cell Cell, f func(int) int) ComputeCell {
	c1 := cell.RegisterListener()

	compute := &Compute{
		sheet:  s,
		value:  f(cell.Value()),
		inputs: []chan int{c1},
		f1:     f,
		cells:  []Cell{cell},
	}

	s.deps[cell] = dependencies{
		f1:      f,
		targets: []Cell{compute},
		deps:    []Cell{cell},
	}

	go compute.listen()

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}

func (s *Sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	c1 := cell.RegisterListener()
	c2 := cell2.RegisterListener()

	compute := &Compute{
		sheet:  s,
		value:  f(cell.Value(), cell2.Value()),
		inputs: []chan int{c1, c2},
		f2:     f,
		cells:  []Cell{cell, cell2},
	}

	s.deps[cell] = dependencies{f2: f, targets: []Cell{compute}, f2Deps: [][]Cell{{cell, cell2}}}
	s.deps[cell2] = dependencies{f2: f, targets: []Cell{compute}, f2Deps: [][]Cell{{cell, cell2}}}

	go compute.listen()

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}
