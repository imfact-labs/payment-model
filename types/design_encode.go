package types

import (
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (de *Design) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)

	return nil
}
