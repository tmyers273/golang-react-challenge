package react

var _ ComputeCell = &Compute{}

type Compute struct {
	id        int
	value     int
	f1        func(int) int
	f2        func(int, int) int
	inputs    []chan int
	listeners []chan int
	cells     []Cell
	callback  func(int)
}

func (c *Compute) SetValue(i int) {
	if c.value == i {
		return
	}

	c.value = i

	for _, ch := range c.listeners {
		ch <- c.value
	}

	if c.callback != nil {
		c.callback(c.value)
	}
}

func (c *Compute) listen() {
	if c.f1 != nil {
		for {
			select {
			case _ = <-c.inputs[0]:
				c.SetValue(c.f1(c.cells[0].Value()))
			}
		}
	} else if c.f2 != nil {
		for {
			select {
			case _ = <-c.inputs[0]:
				c.SetValue(c.f2(c.cells[0].Value(), c.cells[1].Value()))
			case _ = <-c.inputs[1]:
				c.SetValue(c.f2(c.cells[0].Value(), c.cells[1].Value()))
			}
		}
	}
}

func (c *Compute) GetId() int {
	return c.id
}

func (c *Compute) Value() int {
	return c.value
}

func (c *Compute) RegisterListener() chan int {
	ch := make(chan int)
	c.listeners = append(c.listeners, ch)
	return ch
}

func (c *Compute) AddCallback(f func(int)) Canceler {
	c.callback = f

	return &DummyCanceler{}
}

type DummyCanceler struct{}

func (c *DummyCanceler) Cancel() {}
