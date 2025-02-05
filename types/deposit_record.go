package types

import (
	"encoding/json"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var DepositRecordHint = hint.MustNewHint("mitum-payment-deposit-record-v0.0.1")

type DepositRecord struct {
	hint.BaseHinter
	address base.Address
	items   map[string]DepositRecordItem
}

func NewDepositRecord(
	address base.Address,
) DepositRecord {
	items := make(map[string]DepositRecordItem)
	return DepositRecord{
		BaseHinter: hint.NewBaseHinter(DepositRecordHint),
		address:    address,
		items:      items,
	}
}

func NewEmptyDepositRecord() DepositRecord {
	return DepositRecord{
		BaseHinter: hint.NewBaseHinter(DepositRecordHint),
	}
}

func (d DepositRecord) IsValid([]byte) error {
	if err := d.BaseHinter.IsValid(nil); err != nil {
		return err
	}

	if err := util.CheckIsValiders(nil, false,
		d.address,
	); err != nil {
		return err
	}

	for _, itm := range d.items {
		if err := util.CheckIsValiders(nil, false,
			itm,
		); err != nil {
			return err
		}
	}

	return nil
}

func (d DepositRecord) Bytes() []byte {
	var itm []byte
	if d.items != nil {
		b, _ := json.Marshal(d.items)
		itm = valuehash.NewSHA256(b).Bytes()
	} else {
		itm = []byte{}
	}

	return util.ConcatBytesSlice(
		d.address.Bytes(),
		itm,
	)
}

func (d DepositRecord) Address() base.Address {
	return d.address
}

func (d *DepositRecord) SetItem(cid string, am common.Big, ts uint64) {
	d.items[cid] = NewDepositRecordItem(am, ts)
}

func (d DepositRecord) Amount(cid string) *common.Big {
	itm, found := d.items[cid]
	if !found {
		return nil
	}

	return &itm.Amount
}

func (d DepositRecord) TransferredAt(cid string) *uint64 {
	itm, found := d.items[cid]
	if !found {
		return nil
	}

	return &itm.TransferredAt
}

type DepositRecordItem struct {
	Amount        common.Big `bson:"amount" json:"amount"`
	TransferredAt uint64     `bson:"transferred_at" json:"transferred_at"`
}

func NewDepositRecordItem(am common.Big, ts uint64) DepositRecordItem {
	return DepositRecordItem{
		Amount:        am,
		TransferredAt: ts,
	}
}

func (d DepositRecordItem) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		d.Amount,
	); err != nil {
		return err
	}

	return nil
}
