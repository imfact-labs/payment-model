package types

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (d *DepositRecord) unpack(
	enc encoder.Encoder,
	ht hint.Hint,
	addr string,
) error {
	d.BaseHinter = hint.NewBaseHinter(ht)
	address, err := base.DecodeAddress(addr, enc)
	if err != nil {
		return err
	}
	d.address = address

	return nil
}
