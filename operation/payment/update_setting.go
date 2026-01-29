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
	UpdateAccountSettingFactHint = hint.MustNewHint("mitum-payment-update-account-setting-operation-fact-v0.0.1")
	UpdateAccountSettingHint     = hint.MustNewHint("mitum-payment-update-account-setting-operation-v0.0.1")
)

type UpdateAccountSettingFact struct {
	base.BaseFact
	sender        base.Address
	contract      base.Address
	transferLimit common.Big
	startTime     uint64
	endTime       uint64
	duration      uint64
	currency      ctypes.CurrencyID
}

func NewUpdateAccountSettingFact(
	token []byte, sender, contract base.Address,
	transferLimit common.Big, starTime, endTime, duration uint64, currency ctypes.CurrencyID) UpdateAccountSettingFact {
	bf := base.NewBaseFact(UpdateAccountSettingFactHint, token)
	fact := UpdateAccountSettingFact{
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

func (fact UpdateAccountSettingFact) IsValid(b []byte) error {
	if fact.sender.Equal(fact.contract) {
		return common.ErrFactInvalid.Wrap(
			common.ErrSelfTarget.Wrap(errors.Errorf("sender %v is same with contract account", fact.sender)))
	}

	if fact.endTime == 0 {
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

func (fact UpdateAccountSettingFact) Hash() util.Hash {
	return fact.BaseFact.Hash()
}

func (fact UpdateAccountSettingFact) GenerateHash() util.Hash {
	return valuehash.NewSHA256(fact.Bytes())
}

func (fact UpdateAccountSettingFact) Bytes() []byte {
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

func (fact UpdateAccountSettingFact) Token() base.Token {
	return fact.BaseFact.Token()
}

func (fact UpdateAccountSettingFact) Sender() base.Address {
	return fact.sender
}

func (fact UpdateAccountSettingFact) Contract() base.Address {
	return fact.contract
}

func (fact UpdateAccountSettingFact) TransferLimit() common.Big {
	return fact.transferLimit
}

func (fact UpdateAccountSettingFact) StartTime() uint64 {
	return fact.startTime
}

func (fact UpdateAccountSettingFact) EndTime() uint64 {
	return fact.endTime
}

func (fact UpdateAccountSettingFact) Duration() uint64 {
	return fact.duration
}

func (fact UpdateAccountSettingFact) Currency() ctypes.CurrencyID {
	return fact.currency
}

func (fact UpdateAccountSettingFact) Addresses() ([]base.Address, error) {
	return []base.Address{fact.sender}, nil
}

func (fact UpdateAccountSettingFact) FeeBase() map[ctypes.CurrencyID][]common.Big {
	required := make(map[ctypes.CurrencyID][]common.Big)
	required[fact.Currency()] = []common.Big{common.ZeroBig}

	return required
}

func (fact UpdateAccountSettingFact) FeePayer() base.Address {
	return fact.sender
}

func (fact UpdateAccountSettingFact) FeeItemCount() (uint, bool) {
	return extras.ZeroItem, extras.HasNoItem
}

func (fact UpdateAccountSettingFact) FactUser() base.Address {
	return fact.sender
}

func (fact UpdateAccountSettingFact) Signer() base.Address {
	return fact.sender
}

func (fact UpdateAccountSettingFact) ActiveContract() []base.Address {
	return []base.Address{fact.contract}
}

type UpdateAccountSetting struct {
	extras.ExtendedOperation
}

func NewUpdateAccountSetting(fact UpdateAccountSettingFact) (UpdateAccountSetting, error) {
	return UpdateAccountSetting{
		ExtendedOperation: extras.NewExtendedOperation(UpdateAccountSettingHint, fact),
	}, nil
}
