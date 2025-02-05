package deposit

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type TransferFactJSONMarshaler struct {
	base.BaseFactJSONMarshaler
	Sender   base.Address      `json:"sender"`
	Contract base.Address      `json:"contract"`
	Receiver base.Address      `json:"receiver"`
	Amount   common.Big        `json:"amount"`
	Currency ctypes.CurrencyID `json:"currency"`
}

func (fact TransferFact) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(TransferFactJSONMarshaler{
		BaseFactJSONMarshaler: fact.BaseFact.JSONMarshaler(),
		Sender:                fact.sender,
		Contract:              fact.contract,
		Receiver:              fact.receiver,
		Amount:                fact.amount,
		Currency:              fact.currency,
	})
}

type TransferFactJSONUnmarshaler struct {
	base.BaseFactJSONUnmarshaler
	Sender   string     `json:"sender"`
	Contract string     `json:"contract"`
	Receiver string     `json:"receiver"`
	Amount   common.Big `json:"amount"`
	Currency string     `json:"currency"`
}

func (fact *TransferFact) DecodeJSON(b []byte, enc encoder.Encoder) error {
	var u TransferFactJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	fact.BaseFact.SetJSONUnmarshaler(u.BaseFactJSONUnmarshaler)
	fact.amount = u.Amount
	if err := fact.unpack(
		enc, u.Sender, u.Contract, u.Receiver, u.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeJson, *fact)
	}

	return nil
}

func (op Transfer) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(OperationMarshaler{
		BaseOperationJSONMarshaler:           op.BaseOperation.JSONMarshaler(),
		BaseOperationExtensionsJSONMarshaler: op.BaseOperationExtensions.JSONMarshaler(),
	})
}

func (op *Transfer) DecodeJSON(b []byte, enc encoder.Encoder) error {
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
