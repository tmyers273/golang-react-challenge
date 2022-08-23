package react

var _ Reactor = &Sheet{}

func New() *Sheet {
	s := &Sheet{
		ch: make(chan Update),
	}

	//go s.listen()

	return s
}

type Sheet struct {
	Inputs       []*Input
	ComputeCells []*Compute

	ch chan Update
	id int
}

//func (s *Sheet) listen() {
//	for update := range s.ch {
//		for _, c := range s.ComputeCells {
//			if c.id == update.id {
//				c.value = c.f1(update.value)
//				spew.Dump("should call f2 with", update)
//			}
//		}
//	}
//}

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
		dependencies: []int{cell.GetId()},
		value:        f(cell.Value()),
		inputs:       []chan int{c1},
		f1:           f,
	}

	go compute.listen()

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}

func (s *Sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	c1 := cell.RegisterListener()
	c2 := cell2.RegisterListener()

	compute := &Compute{
		dependencies: []int{cell.GetId(), cell.GetId()},
		value:        f(cell.Value(), cell2.Value()),
		inputs:       []chan int{c1, c2},
		f2:           f,
	}

	s.ComputeCells = append(s.ComputeCells, compute)
	return compute
}
