package react

var _ Reactor = &Sheet{}

func New() *Sheet {
	s := &Sheet{
		ch: make(chan Update),
	}

	return s
}

type Sheet struct {
	Inputs       []*Input
	ComputeCells []*Compute

	ch chan Update
	id int
}

type Update struct {
	id    int
	value int
}

func (s *Sheet) CreateInput(i int) InputCell {
	cell := Input{value: i, id: s.id, ch: s.ch}
	s.Inputs = append(s.Inputs, &cell)
	s.id++
	return &cell
}

func (s *Sheet) CreateCompute1(cell Cell, f func(int) int) ComputeCell {
	c1 := cell.RegisterListener()

	compute := &Compute{
		value:  f(cell.Value()),
		inputs: []chan int{c1},
		f1:     f,
		cells:  []Cell{cell},
	}

	go compute.listen()

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}

func (s *Sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	c1 := cell.RegisterListener()
	c2 := cell2.RegisterListener()

	compute := &Compute{
		value:  f(cell.Value(), cell2.Value()),
		inputs: []chan int{c1, c2},
		f2:     f,
		cells:  []Cell{cell, cell2},
	}

	go compute.listen()

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}
