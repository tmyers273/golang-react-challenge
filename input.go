package react

var _ InputCell = &Input{}

type Input struct {
	id          int
	value       int
	ch          chan<- Update
	subscribers []*Compute
}

func (i *Input) GetId() int {
	return i.id
}

func (i *Input) RegisterSubscriber(c *Compute) {
	i.subscribers = append(i.subscribers, c)
}

func (i *Input) Value() int {
	return i.value
}

func (i *Input) SetValue(i2 int) {
	i.value = i2

	i.ch <- Update{
		id:    i.id,
		value: i.value,
	}
}
