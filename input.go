package react

var _ InputCell = &Input{}

type Input struct {
	value int
}

func (i *Input) Value() int {
	return i.value
}

func (i *Input) SetValue(i2 int) {
	i.value = i2
}
