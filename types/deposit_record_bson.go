package types

import (
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (d DepositRecord) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":   d.Hint().String(),
		"address": d.address,
		"items":   d.items,
	})
}

type DepositRecordBSONUnmarshaler struct {
	Hint    string                       `bson:"_hint"`
	Address string                       `bson:"address"`
	Items   map[string]DepositRecordItem `bson:"items"`
}

func (d *DepositRecord) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DepositRecord")

	var u DepositRecordBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	d.items = u.Items

	err = d.unpack(enc, ht, u.Address)
	if err != nil {
		return e.Wrap(err)
	}

	return nil
}
