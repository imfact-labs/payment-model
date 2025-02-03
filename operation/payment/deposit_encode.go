package payment

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *DepositFact) unpack(
	enc encoder.Encoder,
	sa, ca, tl string,
	st, et, dur uint64,
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

	big, err := common.NewBigFromString(tl)
	if err != nil {
		return err
	}

	fact.transferLimit = big
	fact.startTime = st
	fact.endTime = et
	fact.duration = dur

	return nil
}
