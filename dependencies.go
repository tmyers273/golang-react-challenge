package react

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
)

type f1Dependency struct {
	f      func(int) int
	target *compute
	deps   []Cell
}

func (d *f1Dependency) Hash() string {
	return hash(d.deps)
}

func (d *f1Dependency) Target() *compute {
	return d.target
}

func (d *f1Dependency) Calculate() int {
	return d.f(d.deps[0].Value())
}

var _ dependency = &f2Dependency{}

type f2Dependency struct {
	f      func(int, int) int
	target *compute
	deps   []Cell
}

func (d *f2Dependency) Hash() string {
	return hash(d.deps)
}

func (d *f2Dependency) Target() *compute {
	return d.target
}

func (d *f2Dependency) Calculate() int {
	return d.f(d.deps[0].Value(), d.deps[1].Value())
}

// hash returns a hash of the given slice of cells.
// @todo This feels pretty hacky, required a unique id on each cell.
// It seems like there should be some better way to do this...
func hash(in []Cell) string {
	out := make([]byte, len(in)*8)
	for i := 0; i < len(in); i += 8 {
		binary.LittleEndian.PutUint64(out[i:i+8], uint64(in[i].GetId()))
	}

	hashed := md5.Sum(out)
	return hex.EncodeToString(hashed[:])
}
