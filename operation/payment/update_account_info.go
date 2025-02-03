package payment

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"github.com/pkg/errors"
)

var (
	UpdateAccountInfoFactHint = hint.MustNewHint("mitum-payment-update-account-info-fact-operation-fact-v0.0.1")
	UpdateAccountInfoHint     = hint.MustNewHint("mitum-payment-update-account-info-operation-v0.0.1")
)

type UpdateAccountInfoFact struct {
	base.BaseFact
	sender        base.Address
	contract      base.Address
	transferLimit ctypes.Amount
	startTime     uint64
	endTime       uint64
	duration      uint64
	currency      ctypes.CurrencyID
}

func NewUpdateAccountInfoFact(
	token []byte, sender, contract base.Address,
	transferLimit ctypes.Amount, starTime, endTime, duration uint64, currency ctypes.CurrencyID) UpdateAccountInfoFact {
	bf := base.NewBaseFact(UpdateAccountInfoFactHint, token)
	fact := UpdateAccountInfoFact{
		BaseFact:      bf,
		sender:        sender,
		contract:      contract,
		transferLimit: transferLimit,
		startTime:     starTime,
		endTime:       endTime,
		duration:      duration,
		currency:      currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact UpdateAccountInfoFact) IsValid(b []byte) error {
	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if fact.startTime < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Errorf("start time must be bigger than zero"))
	}

	if fact.endTime < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Errorf("end time must be bigger than zero"))
	}

	if fact.duration < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Errorf("duration must be bigger than zero"))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.transferLimit,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact UpdateAccountInfoFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdateAccountInfoFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateAccountInfoFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.transferLimit.Bytes(),
		util.Uint64ToBytes(fact.startTime),
		util.Uint64ToBytes(fact.endTime),
		util.Uint64ToBytes(fact.duration),
		fact.currency.Bytes(),
	)
}

func (fact UpdateAccountInfoFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact UpdateAccountInfoFact) Sender() base.Address {
	return fact.sender
}

func (fact UpdateAccountInfoFact) Contract() base.Address {
	return fact.contract
}

func (fact UpdateAccountInfoFact) TransferLimit() ctypes.Amount {
	return fact.transferLimit
}

func (fact UpdateAccountInfoFact) StartTime() uint64 {
	return fact.startTime
}

func (fact UpdateAccountInfoFact) EndTime() uint64 {
	return fact.endTime
}

func (fact UpdateAccountInfoFact) Duration() uint64 {
	return fact.duration
}

func (fact UpdateAccountInfoFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact UpdateAccountInfoFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender}, nil
}

func (fact UpdateAccountInfoFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact UpdateAccountInfoFact) FeePayer() base.Address {
	return fact.sender
}

func (fact UpdateAccountInfoFact) FactUser() base.Address {
	return fact.sender
}

func (fact UpdateAccountInfoFact) Signer() base.Address {
	return fact.sender
}

func (fact UpdateAccountInfoFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type UpdateAccountInfo struct {
	extras.ExtendedOperation
}

func NewUpdateAccountInfo(fact UpdateAccountInfoFact) (UpdateAccountInfo, error) {
	return UpdateAccountInfo{
		ExtendedOperation: extras.NewExtendedOperation(UpdateAccountInfoHint, fact),
	}, nil
}
