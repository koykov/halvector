package halvector

import (
	"io"

	"github.com/koykov/vector"
)

type Codec struct {
	vector.BaseCodec
}

func (Codec) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (Codec) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
