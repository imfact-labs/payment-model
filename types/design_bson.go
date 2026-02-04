package types

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de Design) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":             de.Hint().String(),
			"transfer_settings": de.settings,
		})
}

type DesignBSONUnmarshaler struct {
	Hint     string   `bson:"_hint"`
	Accounts bson.Raw `bson:"transfer_settings"`
}

func (de *Design) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of Design")

	var u DesignBSONUnmarshaler
	if err := bson.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	accounts := make(map[string]Setting)
	m, err := enc.DecodeMap(u.Accounts)
	if err != nil {
		return e.Wrap(err)
	}
	for k, v := range m {
		ac, ok := v.(Setting)
		if !ok {
			return e.Wrap(errors.Errorf("expected Setting, not %T", v))
		}

		accounts[k] = ac
	}
	de.settings = accounts

	err = de.unpack(enc, ht)
	if err != nil {
		return e.Wrap(err)
	}

	return nil
}
