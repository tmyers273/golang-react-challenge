package react

var _ InputCell = &input{}

type input struct {
	id    int
	value int
	sheet *sheet
}

func (i *input) GetId() int {
	return i.id
}

func (i *input) Value() int {
	return i.value
}

func (i *input) SetValue(i2 int) {
	if i.value == i2 {
		return
	}

	i.value = i2

	fireDeps(i.sheet, i)
}

func fireDeps(sheet *sheet, cell Cell) {
	var old int

	// First we "land" all the new values

	changes := make(map[string]*compute)
	// And keep track of any tuples of referentiality
	// that change. This is used to dedupe the updates.
	if depList, ok := sheet.deps[cell]; ok {
		for _, dep := range depList {
			old = dep.Target().value
			dep.Target().value = dep.Calculate()

			if old != dep.Target().value {
				changes[dep.Hash()] = dep.Target()
			}
		}
	}

	// Then fire deps for anything that has changed. Since we
	// are pulling from a map and have hashed appropriately,
	// this ensures we are only firing for each dependency once.
	for _, c := range changes {
		fireDeps(sheet, c)

		for _, cb := range c.callbacks {
			cb.f(c.value)
		}
	}
}
