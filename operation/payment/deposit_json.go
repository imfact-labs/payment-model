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

type DepositFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender        base.Address  `json:"sender"`
	Contract      base.Address  `json:"contract"`
	Amount        ctypes.Amount `json:"amount"`
	TransferLimit common.Big    `json:"transfer_limit"`
	StartTime     uint64        `json:"start_time"`
	EndTime       uint64        `json:"end_time"`
	Duration      uint64        `json:"duration"`
}

func (fact DepositFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DepositFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Amount:                fact.amount,
		TransferLimit:         fact.transferLimit,
		StartTime:             fact.startTime,
		EndTime:               fact.endTime,
		Duration:              fact.duration,
	})
}

type DepositFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender        string          `json:"sender"`
	Contract      string          `json:"contract"`
	Amount        json.RawMessage `json:"amount"`
	TransferLimit string          `json:"transfer_limit"`
	StartTime     uint64          `json:"start_time"`
	EndTime       uint64          `json:"end_time"`
	Duration      uint64          `json:"duration"`
}

func (fact *DepositFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u DepositFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)

	var amount ctypes.Amount
	err := amount.DecodeJSON(u.Amount, enc)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}
	fact.amount = amount

	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.TransferLimit, u.StartTime, u.EndTime, u.Duration,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

type OperationMarshaler struct {
	common.BaseOperationJSONMarshaler
	extras.BaseOperationExtensionsJSONMarshaler
}

func (op Deposit) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Deposit) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
