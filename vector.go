package halvector

import (
	"github.com/koykov/byteconv"
	"github.com/koykov/vector"
)

const (
	flagSorted = 8
)

type Vector struct {
	vector.Vector
	limit int
}

func NewVector() *Vector {
	vec := &Vector{}
	vec.SetBit(vector.FlagInit, true)
	vec.Helper = helper
	return vec
}

func (vec *Vector) Parse(s []byte) error {
	return vec.parse(s, false)
}

func (vec *Vector) ParseString(s string) error {
	return vec.parse(byteconv.S2B(s), false)
}

func (vec *Vector) ParseCopy(s []byte) error {
	return vec.parse(s, true)
}

func (vec *Vector) ParseCopyString(s string) error {
	return vec.parse(byteconv.S2B(s), true)
}

// SetLimit setups hard of nodes. All entities over the limit will ignore.
// See BenchmarkLimit() for explanation.
func (vec *Vector) SetLimit(limit int) *Vector {
	if limit < 0 {
		limit = 0
	}
	vec.limit = limit
	return vec
}

func (vec *Vector) Reset() {
	vec.Vector.Reset()
	vec.limit = 0
}
