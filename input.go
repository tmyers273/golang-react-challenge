package react

var _ InputCell = &Input{}

type Input struct {
	id        int
	value     int
	ch        chan<- Update
	listeners []chan int
	sheet     *Sheet
}

func (i *Input) RegisterListener() chan int {
	ch := make(chan int)
	i.listeners = append(i.listeners, ch)
	return ch
}

func (i *Input) GetId() int {
	return i.id
}

func (i *Input) Value() int {
	return i.value
}

func (i *Input) SetValue(i2 int) {
	if i.value == i2 {
		return
	}

	i.value = i2

	for _, c := range i.listeners {
		c <- i.value
	}
}
