package react

var _ Reactor = &Sheet{}

func New() *Sheet {
	return &Sheet{}
}

type Sheet struct {
	Inputs       []*Input
	ComputeCells []*ComputeCell
}

func (s *Sheet) CreateInput(i int) InputCell {
	cell := Input{value: i}
	s.Inputs = append(s.Inputs, &cell)
	return &cell
}

func (s *Sheet) CreateCompute1(cell Cell, f func(int) int) ComputeCell {
	//TODO implement me
	panic("implement me")
}

func (s *Sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	//TODO implement me
	panic("implement me")
}
