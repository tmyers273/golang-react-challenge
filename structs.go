package react

type callback struct {
	id int
	f  func(int)
}

func newCallbackRemover(id int, c *compute) Canceler {
	return &callbackRemover{
		id:   id,
		cell: c,
	}
}

type callbackRemover struct {
	id   int
	cell *compute
}

func (c *callbackRemover) Cancel() {
	c.cell.removeCallback(c.id)
}
