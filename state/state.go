package state

import (
	"fmt"
	"strings"

	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/payment-model/types"
	"github.com/pkg/errors"
)

var (
	DesignStateValueHint  = hint.MustNewHint("mitum-payment-design-state-value-v0.0.1")
	PaymentStateKeyPrefix = "payment"
	DesignStateKeySuffix  = "design"
)

func PaymentStateKey(addr string) string {
	return fmt.Sprintf("%s:%s", PaymentStateKeyPrefix, addr)
}

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (sv DesignStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := sv.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv DesignStateValue) HashBytes() []byte {
	return sv.Design.Bytes()
}

func GetDesignFromState(st base.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, errors.Errorf("state value is nil")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("expected DesignStateValue but %T", v)
	}

	return d.Design, nil
}

func IsDesignStateKey(key string) bool {
	return strings.HasPrefix(key, PaymentStateKeyPrefix) && strings.HasSuffix(key, DesignStateKeySuffix)
}

func DesignStateKey(addr string) string {
	return fmt.Sprintf("%s:%s", PaymentStateKey(addr), DesignStateKeySuffix)
}

var (
	DepositRecordStateValueHint = hint.MustNewHint("mitum-payment-deposit-record-state-value-v0.0.1")
	DepositRecordStateKeySuffix = "depositrecord"
)

type DepositRecordStateValue struct {
	hint.BaseHinter
	Record types.DepositRecord
}

func NewDepositRecordStateValue(record types.DepositRecord) DepositRecordStateValue {
	return DepositRecordStateValue{
		BaseHinter: hint.NewBaseHinter(DepositRecordStateValueHint),
		Record:     record,
	}
}

func (sv DepositRecordStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv DepositRecordStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DepositRecordStateValue")

	if err := sv.BaseHinter.IsValid(DepositRecordStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv DepositRecordStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice(sv.Record.Bytes())
}

func GetDepositRecordFromState(st base.State) (*types.DepositRecord, error) {
	v := st.Value()
	if v == nil {
		return nil, errors.Errorf("state value is nil")
	}

	isv, ok := v.(DepositRecordStateValue)
	if !ok {
		return nil, errors.Errorf("expected DepositRecordStateValue but, %T", v)
	}

	return &isv.Record, nil
}

func IsDepositRecordStateKey(key string) bool {
	return strings.HasPrefix(key, PaymentStateKeyPrefix) && strings.HasSuffix(key, DepositRecordStateKeySuffix)
}

func DepositRecordStateKey(addr string, acAddr string) string {
	return fmt.Sprintf("%s:%s:%s", PaymentStateKey(addr), acAddr, DepositRecordStateKeySuffix)
}
