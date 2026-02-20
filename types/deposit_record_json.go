package types

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

type DepositRecordJSONMarshaler struct {
	hint.BaseHinter
	Address base.Address                 `json:"address"`
	Items   map[string]DepositRecordItem `json:"items"`
}

func (d DepositRecord) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(DepositRecordJSONMarshaler{
		BaseHinter: d.BaseHinter,
		Address:    d.address,
		Items:      d.items,
	})
}

type DepositRecordJSONUnmarshaler struct {
	Hint    hint.Hint                    `json:"_hint"`
	Address string                       `json:"address"`
	Items   map[string]DepositRecordItem `json:"items"`
}

func (d *DepositRecord) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of DepositRecord")

	var u DepositRecordJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	d.items = u.Items
	err := d.unpack(enc, u.Hint, u.Address)
	if err != nil {
		return e.Wrap(err)
	}

	return nil
}
