package deposit

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
	WithdrawFactHint = hint.MustNewHint("mitum-payment-withdraw-operation-fact-v0.0.1")
	WithdrawHint     = hint.MustNewHint("mitum-payment-withdraw-operation-v0.0.1")
)

type WithdrawFact struct {
	base.BaseFact
	sender          base.Address
	contract        base.Address
	depositCurrency ctypes.CurrencyID
	currency        ctypes.CurrencyID
}

func NewWithdrawFact(
	token []byte, sender, contract base.Address, dCurrency, currency ctypes.CurrencyID) WithdrawFact {
	bf := base.NewBaseFact(WithdrawFactHint, token)
	fact := WithdrawFact{
		BaseFact:        bf,
		sender:          sender,
		contract:        contract,
		depositCurrency: dCurrency,
		currency:        currency,
	}

	fact.SetHash(fact.GenerateHash())
	return fact
}

func (fact WithdrawFact) IsValid(b []byte) error {
	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if err := util.CheckIsValiders(nil, false,
		fact.BaseHinter,
		fact.sender,
		fact.contract,
		fact.depositCurrency,
		fact.currency,
	); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	if err := common.IsValidOperationFact(fact, b); err != nil {
		return common.ErrFactInvalid.Wrap(err)
	}

	return nil
}

func (fact WithdrawFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact WithdrawFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact WithdrawFact) Bytes() []byte {
	return util.ConcatBytesSlice(
		fact.Token(),
		fact.sender.Bytes(),
		fact.contract.Bytes(),
		fact.depositCurrency.Bytes(),
		fact.currency.Bytes(),
	)
}

func (fact WithdrawFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact WithdrawFact) Sender() base.Address {
	return fact.sender
}

func (fact WithdrawFact) Contract() base.Address {
	return fact.contract
}

func (fact WithdrawFact) DepositCurrency() ctypes.CurrencyID {
	return fact.depositCurrency
}

func (fact WithdrawFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact WithdrawFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender}, nil
}

func (fact WithdrawFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact WithdrawFact) FeePayer() base.Address {
	return fact.sender
}

func (fact WithdrawFact) FactUser() base.Address {
	return fact.sender
}

func (fact WithdrawFact) Signer() base.Address {
	return fact.sender
}

func (fact WithdrawFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type Withdraw struct {
	extras.ExtendedOperation
}

func NewWithdraw(fact WithdrawFact) (Withdraw, error) {
	return Withdraw{
		ExtendedOperation: extras.NewExtendedOperation(WithdrawHint, fact),
	}, nil
}
