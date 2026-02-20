package payment

import (
	"github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util/encoder"
)

func (fact *UpdateAccountSettingFact) unpack(
	enc encoder.Encoder,
	sa, ca string,
	st, et, dur uint64,
	cid string,
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
	fact.currency = types.CurrencyID(cid)

	return nil
}
