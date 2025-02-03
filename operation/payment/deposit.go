package payment

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
	DepositFactHint = hint.MustNewHint("mitum-payment-deposit-operation-fact-v0.0.1")
	DepositHint     = hint.MustNewHint("mitum-payment-deposit-operation-v0.0.1")
)

type DepositFact struct {
	base.BaseFact
	sender        base.Address
	contract      base.Address
	amount        ctypes.Amount
	transferLimit common.Big
	startTime     uint64
	endTime       uint64
	duration      uint64
}

func NewDepositFact(
	token []byte,
	sender, contract base.Address,
	amount ctypes.Amount,
	transferLimit common.Big,
	startTime, endTime, duration uint64,

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
	)
}

func (fact DepositFact) IsValid(b []byte) error {
	if err := fact.BaseHinter.IsValid(nil); err != nil {
		return common.ErrFactInvalid.Wrap(err)
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
		fact.sender,
		fact.contract,
		fact.amount,
		fact.transferLimit,
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

func (fact DepositFact) Amount() ctypes.Amount {
	return fact.amount
}

func (fact DepositFact) TransferLimit() common.Big {
	return fact.transferLimit
}

func (fact DepositFact) Signer() base.Address {
	return fact.sender
}

func (fact DepositFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.Sender()}, nil
}

func (fact DepositFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Amount().Currency()] = []common.Big{fact.Amount().Big()}

	return required
}

func (fact DepositFact) FeePayer() base.Address {
	return fact.sender
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
