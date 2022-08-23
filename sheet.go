package react

var _ Reactor = &sheet{}

func New() *sheet {
	s := &sheet{
		deps: make(map[Cell][]dependency),
	}

	return s
}

type sheet struct {
	Inputs       []*input
	ComputeCells []*compute

	deps map[Cell][]dependency

	id int
}

func (s *sheet) CreateInput(i int) InputCell {
	cell := input{value: i, id: s.id, sheet: s}
	s.Inputs = append(s.Inputs, &cell)
	s.id++
	return &cell
}

func (s *sheet) CreateCompute1(cell Cell, f func(int) int) ComputeCell {
	computedCell := &compute{
		sheet: s,
		value: f(cell.Value()),
	}

	dep := &f1Dependency{
		f:      f,
		target: computedCell,
		deps:   []Cell{cell},
	}

	s.deps[cell] = append(s.deps[cell], dep)

	s.ComputeCells = append(s.ComputeCells, computedCell)
	return computedCell
}

func (s *sheet) CreateCompute2(cell Cell, cell2 Cell, f func(int, int) int) ComputeCell {
	computedCell := &compute{
		sheet: s,
		value: f(cell.Value(), cell2.Value()),
	}

	dep := &f2Dependency{
		f:      f,
		target: computedCell,
		deps:   []Cell{cell, cell2},
	}

	s.deps[cell] = append(s.deps[cell], dep)
	s.deps[cell2] = append(s.deps[cell2], dep)

	s.ComputeCells = append(s.ComputeCells, computedCell)
	return computedCell
}
