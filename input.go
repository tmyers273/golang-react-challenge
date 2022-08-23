package react

import "fmt"

var _ InputCell = &Input{}

type Input struct {
	id          int
	value       int
	ch          chan<- Update
	subscribers []*Compute
	listeners   []chan int
}

func (i *Input) RegisterListener() chan int {
	ch := make(chan int)
	i.listeners = append(i.listeners, ch)
	return ch
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

	fmt.Printf("pushing %d on %d listeners\n", i.value, len(i.listeners))
	for _, c := range i.listeners {
		c <- i.value
	}
	fmt.Printf("done pushing")
	//i.ch <- Update{
	//	id:    i.id,
	//	value: i.value,
	//}
}
