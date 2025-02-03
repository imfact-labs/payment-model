package payment

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type UpdateAccountInfoFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender        base.Address      `json:"sender"`
	Contract      base.Address      `json:"contract"`
	TransferLimit ctypes.Amount     `json:"transfer_limit"`
	StartTime     uint64            `json:"start_time"`
	EndTime       uint64            `json:"end_time"`
	Duration      uint64            `json:"duration"`
	Currency      ctypes.CurrencyID `json:"currency"`
}

func (fact UpdateAccountInfoFact) MarshalJSON() ([]byte, error) {
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
	Sender        string          `json:"sender"`
	Contract      string          `json:"contract"`
	TransferLimit json.RawMessage `json:"transfer_limit"`
	StartTime     uint64          `json:"start_time"`
	EndTime       uint64          `json:"end_time"`
	Duration      uint64          `json:"duration"`
	Currency      string          `json:"currency"`
}

func (fact *UpdateAccountInfoFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u UpdateAccountInfoFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	var transferLimit ctypes.Amount
	err := transferLimit.DecodeJSON(u.TransferLimit, enc)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}
	fact.transferLimit = transferLimit

	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.StartTime, u.EndTime, u.Duration, u.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op UpdateAccountInfo) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *UpdateAccountInfo) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
