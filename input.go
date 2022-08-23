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

	if m, ok := i.sheet.deps[i]; ok {
		if m.f1 != nil {
			for i, target := range m.targets {
				target.SetValue(m.f1(m.deps[i].Value()))
			}
		} else {
			for i, target := range m.targets {
				target.SetValue(m.f2(m.f2Deps[i][0].Value(), m.f2Deps[i][1].Value()))
			}
		}
	}
}
