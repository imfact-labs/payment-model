package payment

import (
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

func (fact *TransferFact) unpack(
	enc encoder.Encoder,
	sa, ca, ra string,
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

	switch receiver, err := base.DecodeAddress(ra, enc); {
	case err != nil:
		return err
	default:
		fact.receiver = receiver
	}

	return nil
}
