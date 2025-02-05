package deposit

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type WithdrawFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender          base.Address      `json:"sender"`
	Contract        base.Address      `json:"contract"`
	DepositCurrency ctypes.CurrencyID `json:"deposit_currency"`
	Currency        ctypes.CurrencyID `json:"currency"`
}

func (fact WithdrawFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(WithdrawFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		DepositCurrency:       fact.depositCurrency,
		Currency:              fact.currency,
	})
}

type WithdrawFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender          string `json:"sender"`
	Contract        string `json:"contract"`
	DepositCurrency string `json:"deposit_currency"`
	Currency        string `json:"currency"`
}

func (fact *WithdrawFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u WithdrawFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.DepositCurrency, u.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op Withdraw) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Withdraw) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeJSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
