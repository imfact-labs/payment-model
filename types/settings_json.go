package types

import (
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
)

type SettingJSONMarshaler struct {
	hint.BaseHinter
	Address base.Address           `json:"address"`
	Items   map[string]SettingItem `json:"items"`
}

func (s Setting) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(SettingJSONMarshaler{
		BaseHinter: s.BaseHinter,
		Address:    s.address,
		Items:      s.items,
	})
}

type SettingJSONUnmarshaler struct {
	Hint    hint.Hint              `json:"_hint"`
	Address string                 `json:"address"`
	Items   map[string]SettingItem `json:"items"`
}

func (s *Setting) DecodeJSON(b []byte, enc encoder.Encoder) error {
	e := util.StringError("failed to decode json of AccountInfo")

	var u SettingJSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	s.items = u.Items

	err := s.unpack(enc, u.Hint, u.Address)
	if err != nil {
		return e.Wrap(err)
	}

	return nil
}
