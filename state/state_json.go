package state

import (
	"encoding/json"

	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignStateValueJSONMarshaler struct {
	hint.BaseHinter
	Design types.Design `json:"design"`
}

func (sv DesignStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		DesignStateValueJSONMarshaler(sv),
	)
}

type DesignStateValueJSONUnmarshaler struct {
	Hint   hint.Hint       `json:"_hint"`
	Design json.RawMessage `json:"design"`
}

func (sv *DesignStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DesignStateValue")

	var u DesignStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	sv.BaseHinter = hint.NewBaseHinter(u.Hint)

	var sd types.Design
	if err := sd.DecodeJSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Design = sd

	return nil
}

type DepositRecordStateValueJSONMarshaler struct {
	hint.BaseHinter
	DepositRecord types.DepositRecord `json:"deposit_info"`
}

func (sv DepositRecordStateValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(
		DepositRecordStateValueJSONMarshaler(sv),
	)
}

type DepositRecordStateValueJSONUnmarshaler struct {
	Hint          hint.Hint       `json:"_hint"`
	DepositRecord json.RawMessage `json:"deposit_info"`
}

func (sv *DepositRecordStateValue) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DepositRecordStateValue")

	var u DepositRecordStateValueJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	sv.BaseHinter = hint.NewBaseHinter(u.Hint)
	var depositInfo types.DepositRecord
	if err := depositInfo.DecodeJSON(u.DepositRecord, enc); err != nil {
		return e.Wrap(err)
	}
	sv.DepositRecord = depositInfo

	return nil
}
