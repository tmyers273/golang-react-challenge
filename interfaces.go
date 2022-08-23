package react

type Reactor interface {
	CreateInput(int) InputCell
	CreateCompute1(Cell, func(int) int) ComputeCell
	CreateCompute2(Cell, Cell, func(int, int) int) ComputeCell
}

type Cell interface {
	Value() int
	GetId() int // new
}

type InputCell interface {
	Cell
	SetValue(int)
}

type ComputeCell interface {
	Cell
	AddCallback(func(int)) Canceler
}

type Canceler interface {
	Cancel()
}

type dependency interface {
	Calculate() int
	Target() *compute
	Hash() string
}

var _ dependency = &f1Dependency{}
