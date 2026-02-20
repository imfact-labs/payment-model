package payment

import (
	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/currency-model/operation/extras"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
)

type UpdateAccountInfoFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender        base.Address      `json:"sender"`
	Contract      base.Address      `json:"contract"`
	TransferLimit common.Big        `json:"transfer_limit"`
	StartTime     uint64            `json:"start_time"`
	EndTime       uint64            `json:"end_time"`
	Duration      uint64            `json:"duration"`
	Currency      ctypes.CurrencyID `json:"currency"`
}

func (fact UpdateAccountSettingFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(UpdateAccountInfoFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		TransferLimit:         fact.transferLimit,
		StartTime:             fact.startTime,
		EndTime:               fact.endTime,
		Duration:              fact.duration,
		Currency:              fact.currency,
	})
}

type UpdateAccountInfoFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender        string     `json:"sender"`
	Contract      string     `json:"contract"`
	TransferLimit common.Big `json:"transfer_limit"`
	StartTime     uint64     `json:"start_time"`
	EndTime       uint64     `json:"end_time"`
	Duration      uint64     `json:"duration"`
	Currency      string     `json:"currency"`
}

func (fact *UpdateAccountSettingFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u UpdateAccountInfoFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)
	fact.transferLimit = u.TransferLimit

	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.StartTime, u.EndTime, u.Duration, u.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op UpdateAccountSetting) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *UpdateAccountSetting) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
