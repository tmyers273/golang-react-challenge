package react

import (
	"sync"
)

var _ ComputeCell = &Compute{}

type Compute struct {
	id         int
	value      int
	f1         func(int) int
	f2         func(int, int) int
	inputs     []chan int
	listeners  []chan int
	cells      []Cell
	callbacks  []Callback
	mu         sync.Mutex
	callbackId int
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
		f.f(c.value)
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
	c.callbacks = append(c.callbacks, Callback{
		f:  f,
		id: c.callbackId,
	})
	c.callbackId++

	return NewCallbackRemover(len(c.callbacks)-1, c)
}

func (c *Compute) remove(id int) {
	c.mu.Lock()
	for i := 0; i < len(c.callbacks)-1; i++ {
		if c.callbacks[i].id == id {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
			break
		}
	}
	c.mu.Unlock()
}

type Callback struct {
	id int
	f  func(int)
}

func NewCallbackRemover(id int, c *Compute) Canceler {
	return &CallbackRemover{id: id, c: c}
}

type CallbackRemover struct {
	id int
	c  *Compute
}

func (c *CallbackRemover) Cancel() {
	c.c.remove(c.id)
}
