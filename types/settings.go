package types

import (
	"encoding/json"

	"github.com/imfact-labs/currency-model/common"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/valuehash"
)

var SettingHint = hint.MustNewHint("mitum-payment-setting-v0.0.1")

type Setting struct {
	hint.BaseHinter
	address base.Address
	items   map[string]SettingItem
}

func NewSettings(
	address base.Address,
) Setting {
	items := make(map[string]SettingItem)
	return Setting{
		BaseHinter: hint.NewBaseHinter(SettingHint),
		address:    address,
		items:      items,
	}
}

func (s Setting) IsValid([]byte) error {
	if err := s.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		s.address,
	); err != nil {
		return err
	}

	for _, v := range s.items {
		if err := util.CheckIsValiders(nil, false,
			v,
		); err != nil {
			return err
		}
	}

	return nil
}

func (s Setting) Bytes() []byte {
	var itm []byte
	if s.items != nil {
		b, _ := json.Marshal(s.items)
		itm = valuehash.NewSHA256(b).Bytes()
	} else {
		itm = []byte{}
	}

	return util.ConcatBytesSlice(
		s.address.Bytes(),
		itm,
	)
}

func (s Setting) Address() base.Address {
	return s.address
}

func (s Setting) Items() map[string]SettingItem {
	return s.items
}

func (s *Setting) SetItem(cid string, tLimit common.Big, startTime, endTime, duration uint64) {
	s.items[cid] = NewSettingItem(tLimit, startTime, endTime, duration)
}

func (s Setting) TransferLimit(cid string) *common.Big {
	itm, found := s.items[cid]
	if !found {
		return nil
	}

	return &itm.TransferLimit
}

func (s Setting) PeriodTime(cid string) *[3]uint64 {
	pt, found := s.items[cid]
	if !found {
		return nil
	}
	pTime := [3]uint64{pt.StartTime, pt.EndTime, pt.Duration}

	return &pTime
}

func (s *Setting) Remove(cid string) error {
	_, found := s.items[cid]
	if !found {
		return nil
	}
	delete(s.items, cid)

	return nil
}

type SettingItem struct {
	TransferLimit common.Big `bson:"transfer_limit" json:"transfer_limit"`
	StartTime     uint64     `bson:"start_time" json:"start_time"`
	EndTime       uint64     `bson:"end_time" json:"end_time"`
	Duration      uint64     `bson:"duration" json:"duration"`
}

func NewSettingItem(tL common.Big, st, et, dur uint64) SettingItem {
	return SettingItem{
		TransferLimit: tL,
		StartTime:     st,
		EndTime:       et,
		Duration:      dur,
	}
}

func (t SettingItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		t.TransferLimit,
	); err != nil {
		return err
	}
	if t.StartTime < 1 || t.EndTime < 1 || t.Duration < 1 {
		return common.ErrFactInvalid.Wrap(common.ErrValueInvalid.Errorf("time data must be greater than zero"))
	}

	return nil
}
