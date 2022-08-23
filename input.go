package react

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
)

var _ InputCell = &Input{}

type Input struct {
	id        int
	value     int
	ch        chan<- Update
	listeners []chan int
	sheet     *Sheet
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

	fireDeps(i.sheet, i)
}

func hash(in []Cell) string {
	out := make([]byte, len(in)*8)
	for i := 0; i < len(in); i += 8 {
		binary.LittleEndian.PutUint64(out[i:i+8], uint64(in[i].GetId()))
	}

	hashed := md5.Sum(out)
	return hex.EncodeToString(hashed[:])
}

func fireDeps(sheet *Sheet, c Cell) {
	// First "land" all the new values
	var old int
	changes := make(map[string]*Compute)
	if m, ok := sheet.deps[c]; ok {
		for k := 0; k < len(m); k++ {
			if m[k].f1 != nil {
				old = m[k].target.value
				m[k].target.value = m[k].f1(m[k].f1Deps[0].Value())

				if old != m[k].target.value {
					changes[hash(m[k].f1Deps)] = m[k].target
				}
			} else {
				old = m[k].target.value
				m[k].target.value = m[k].f2(m[k].f2Deps[0].Value(), m[k].f2Deps[1].Value())
				if old != m[k].target.value {
					changes[hash(m[k].f1Deps)] = m[k].target
				}
			}
		}
	}

	// Then fire deps for anything that has changed
	for _, c := range changes {

		fireDeps(sheet, c)

		for _, cb := range c.callbacks {
			cb.f(c.value)
		}
	}
}
