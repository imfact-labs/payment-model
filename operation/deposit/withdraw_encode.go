package deposit

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *WithdrawFact) unpack(
	enc encoder.Encoder,
	sa, ca string,
	dcid, cid string,
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

	fact.depositCurrency = types.CurrencyID(dcid)
	fact.currency = types.CurrencyID(cid)

	return nil
}
