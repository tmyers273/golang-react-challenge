package react

var _ ComputeCell = &Compute{}

type Compute struct {
	id           int
	dependencies []int
	value        int
	f1           func(int) int
	f2           func(int, int) int
}

func (c Compute) GetId() int {
	return c.id
}

func (c Compute) Value() int {
	return c.value
}

func (c Compute) AddCallback(f func(int)) Canceler {
	//TODO implement me
	panic("implement me")
}
