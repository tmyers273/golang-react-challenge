package react

import (
	"sync"
)

var _ ComputeCell = &Compute{}

type Compute struct {
	id        int
	value     int
	f1        func(int) int
	f2        func(int, int) int
	inputs    []chan int
	listeners []chan int
	cells     []Cell
	callbacks []func(int)
	mu        sync.Mutex
}

func (c *Compute) SetValue(i int) {
	if c.value == i {
		return
	}

	c.value = i

	for _, ch := range c.listeners {
		ch <- c.value
	}

	for _, f := range c.callbacks {
		f(c.value)
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
	c.callbacks = append(c.callbacks, f)

	return NewCallbackRemover(len(c.callbacks)-1, c)
}

func (c *Compute) remove(index int) {
	c.mu.Lock()
	c.callbacks = append(c.callbacks[:index], c.callbacks[index+1:]...)
	c.mu.Unlock()
}

func NewCallbackRemover(index int, c *Compute) Canceler {
	return &CallbackRemover{index: index, c: c}
}

type CallbackRemover struct {
	index int
	c     *Compute
}

func (c *CallbackRemover) Cancel() {
	c.c.remove(c.index)
}
