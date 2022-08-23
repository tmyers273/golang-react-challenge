package react

var _ ComputeCell = &Compute{}

type Compute struct {
	value int
}

func (c Compute) Value() int {
	//TODO implement me
	panic("implement me")
}

func (c Compute) AddCallback(f func(int)) Canceler {
	//TODO implement me
	panic("implement me")
}
