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
	DepositFactHint = hint.MustNewHint("mitum-payment-deposit-operation-fact-v0.0.1")
	DepositHint     = hint.MustNewHint("mitum-payment-deposit-operation-v0.0.1")
)

type DepositFact struct {
	base.BaseFact
	sender        base.Address
	contract      base.Address
	amount        common.Big
	transferLimit common.Big
	startTime     uint64
	endTime       uint64
	duration      uint64
	currency      ctypes.CurrencyID
}

func NewDepositFact(
	token []byte,
	sender, contract base.Address,
	amount, transferLimit common.Big,
	startTime, endTime, duration uint64, currency ctypes.CurrencyID,

) DepositFact {
	bf := base.NewBaseFact(DepositFactHint, token)
	fact := DepositFact{
		BaseFact:      bf,
		sender:        sender,
		contract:      contract,
		amount:        amount,
		transferLimit: transferLimit,
		startTime:     startTime,
		endTime:       endTime,
		duration:      duration,
		currency:      currency,
	}
	fact.SetHash(fact.GenerateHash())

	return fact
}

func (fact DepositFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact DepositFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact DepositFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.amount.Bytes(),
		fact.transferLimit.Bytes(),
		util.Uint64ToBytes(fact.startTime),
		util.Uint64ToBytes(fact.endTime),
		util.Uint64ToBytes(fact.duration),
		fact.currency.Bytes(),
	)
}

func (fact DepositFact) IsValid(b []byte) error {
	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if fact.Amount().IsZero() {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Errorf("amount cannot be zero"))
	} else if fact.endTime == 0 {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Errorf("end time cannot be zero"))
	} else if fact.startTime >= fact.endTime {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Errorf("start time cannot be greater than end time or equal with end time"))
	} else if fact.duration > (fact.endTime - fact.startTime) {
		return common.ErrFactInvalid.Wrap(
			common.ErrValueInvalid.Errorf("duration cannot be greater than the difference between start and end time"))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.amount,
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

func (fact DepositFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact DepositFact) Sender() base.Address {
	return fact.sender
}

func (fact DepositFact) Contract() base.Address {
	return fact.contract
}

func (fact DepositFact) Amount() common.Big {
	return fact.amount
}

func (fact DepositFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact DepositFact) TransferLimit() common.Big {
	return fact.transferLimit
}

func (fact DepositFact) StartTime() uint64 {
	return fact.startTime
}

func (fact DepositFact) EndTime() uint64 {
	return fact.endTime
}

func (fact DepositFact) Duration() uint64 {
	return fact.duration
}

func (fact DepositFact) Signer() base.Address {
	return fact.sender
}

func (fact DepositFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.Sender()}, nil
}

func (fact DepositFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{fact.Amount()}

	return required
}

func (fact DepositFact) FeePayer() base.Address {
	return fact.sender
}

func (fact DepositFact) FeeItemCount() (uint, bool) {
	return extras.ZeroItem, extras.HasNoItem
}

func (fact DepositFact) FactUser() base.Address {
	return fact.sender
}

func (fact DepositFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type Deposit struct {
	extras.ExtendedOperation
}

func NewDeposit(fact base.Fact) (Deposit, error) {
	return Deposit{
		ExtendedOperation: extras.NewExtendedOperation(DepositHint, fact),
	}, nil
}
