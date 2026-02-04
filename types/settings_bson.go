package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s Setting) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":   s.Hint().String(),
		"address": s.address,
		"items":   s.items,
	})
}

type SettingBSONUnmarshaler struct {
	Hint    string                 `bson:"_hint"`
	Address string                 `bson:"address"`
	Items   map[string]SettingItem `bson:"items"`
}

func (s *Setting) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Setting")

	var u SettingBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	s.items = u.Items

	err = s.unpack(enc, ht, u.Address)
	if err != nil {
		return e.Wrap(err)
	}

	return nil
}
