package deposit

import (
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *DepositFact) unpack(
	enc encoder.Encoder,
	sa, ca string,
	st, et, dur uint64,
	ci string,
) error {
	switch sender, err := base.DecodeAddress(sa, enc); {
	case err != nil:
		return err
	default:
		fact.sender = sender
	}

	switch contract, err := base.DecodeAddress(ca, enc); {
	case err != nil:
		return err
	default:
		fact.contract = contract
	}

	fact.startTime = st
	fact.endTime = et
	fact.duration = dur
	fact.currency = ctypes.CurrencyID(ci)

	return nil
}
