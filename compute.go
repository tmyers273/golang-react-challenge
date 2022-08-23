package react

import (
	"sync"
)

var _ ComputeCell = &compute{}

type compute struct {
	// id is guaranteed unique for all compute cells
	id    int
	value int
	sheet *sheet

	// callbacks is a slice of callback functions, with
	// unique ids. The id is used to remove the callback.
	//
	// A map based (map[int]func(int)) may make more sense
	// here if there is expected to be a lot of churn in
	// the callbacks list.
	callbacks []callback
	mu        sync.Mutex

	// Holds the value of the next callback id. Controlled
	// via the mutex to ensure that the id is unique.
	nextCallbackId int
}

func (c *compute) GetId() int {
	return c.id
}

func (c *compute) Value() int {
	return c.value
}

func (c *compute) AddCallback(f func(int)) Canceler {
	c.mu.Lock()
	c.callbacks = append(c.callbacks, callback{
		f:  f,
		id: c.nextCallbackId,
	})
	c.nextCallbackId++
	c.mu.Unlock()

	return newCallbackRemover(len(c.callbacks)-1, c)
}

func (c *compute) removeCallback(id int) {
	c.mu.Lock()

	for i := 0; i < len(c.callbacks); i++ {
		if c.callbacks[i].id == id {
			c.callbacks = append(c.callbacks[:i], c.callbacks[i+1:]...)
			break
		}
	}

	c.mu.Unlock()
}
