package deposit

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var (
	TransferFactHint = hint.MustNewHint("mitum-payment-transfer-operation-fact-v0.0.1")
	TransferHint     = hint.MustNewHint("mitum-payment-transfer-operation-v0.0.1")
)

type TransferFact struct {
	base.BaseFact
	sender   base.Address
	contract base.Address
	receiver base.Address
	amount   common.Big
	currency ctypes.CurrencyID
}

func NewTransferFact(
	token []byte,
	sender, contract, receiver base.Address,
	amount common.Big, currency ctypes.CurrencyID,
) TransferFact {
	bf := base.NewBaseFact(TransferFactHint, token)
	fact := TransferFact{
		BaseFact: bf,
		sender:   sender,
		contract: contract,
		receiver: receiver,
		amount:   amount,
		currency: currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact TransferFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact TransferFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact TransferFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.receiver.Bytes(),
		fact.amount.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact TransferFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := util.CheckIsValiders(nil, false,
		fact.sender,
		fact.contract,
		fact.receiver,
		fact.amount,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact TransferFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact TransferFact) Sender() base.Address {
	return fact.sender
}

func (fact TransferFact) Contract() base.Address {
	return fact.contract
}

func (fact TransferFact) Receiver() base.Address {
	return fact.receiver
}

func (fact TransferFact) Amount() common.Big {
	return fact.amount
}

func (fact TransferFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact TransferFact) Signer() base.Address {
	return fact.sender
}

func (fact TransferFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.Sender()}, nil
}

func (fact TransferFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{fact.Amount()}

	return required
}

func (fact TransferFact) FeePayer() base.Address {
	return fact.sender
}

func (fact TransferFact) FactUser() base.Address {
	return fact.sender
}

func (fact TransferFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type Transfer struct {
	extras.ExtendedOperation
}

func NewTransfer(fact base.Fact) (Transfer, error) {
	return Transfer{
		ExtendedOperation: extras.NewExtendedOperation(TransferHint, fact),
	}, nil
}
