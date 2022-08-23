package react

var _ ComputeCell = &Compute{}

type Compute struct {
	id           int
	dependencies []int
	value        int
	f1           func(int) int
	f2           func(int, int) int
	inputs       []chan int
	listeners    []chan int
}

func (c *Compute) listen() {
	if c.f1 != nil {
		for {
			select {
			case v := <-c.inputs[0]:
				c.value = c.f1(v)
			}
		}
	} else if c.f2 != nil {
		for {
			select {
			case v := <-c.inputs[0]:
				c.value = c.f2(v, 0)
			case v := <-c.inputs[1]:
				c.value = c.f2(0, v)
			}
		}
	}
}

func (c Compute) GetId() int {
	return c.id
}

func (c Compute) Value() int {
	return c.value
}

func (c Compute) RegisterListener() chan int {
	ch := make(chan int)
	c.listeners = append(c.listeners, ch)
	return ch
}

func (c Compute) AddCallback(f func(int)) Canceler {
	//TODO implement me
	panic("implement me")
}
