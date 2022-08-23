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
	sheet      *Sheet
}

func (c *Compute) SetValue(i int) {
	if c.value == i {
		return
	}

	c.value = i

	fireDeps(c.sheet, c)
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
	c.mu.Lock()
	c.callbacks = append(c.callbacks, Callback{
		f:  f,
		id: c.callbackId,
	})
	c.callbackId++
	c.mu.Unlock()

	return NewCallbackRemover(len(c.callbacks)-1, c)
}

func (c *Compute) remove(id int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := 0; i < len(c.callbacks); i++ {
		if c.callbacks[i].id == id {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
			return
		}
	}
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
