package types

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

func (s *Setting) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	addr string,
) error {
	s.BaseHinter = hint.NewBaseHinter(ht)
	address, err := base.DecodeAddress(addr, enc)
	if err != nil {
		return err
	}
	s.address = address

	return nil
}
