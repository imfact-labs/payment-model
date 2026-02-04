package state

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (sv DesignStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":  sv.Hint().String(),
			"design": sv.Design,
		},
	)
}

type DesignStateValueBSONUnmarshaler struct {
	Hint   string   `bson:"_hint"`
	Design bson.Raw `bson:"design"`
}

func (sv *DesignStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DesignStateValue")

	var u DesignStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	sv.BaseHinter = hint.NewBaseHinter(ht)

	var sd types.Design
	if err := sd.DecodeBSON(u.Design, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Design = sd

	return nil
}

func (sv DepositRecordStateValue) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":          sv.Hint().String(),
			"deposit_record": sv.Record,
		},
	)
}

type DepositRecordStateValueBSONUnmarshaler struct {
	Hint          string   `bson:"_hint"`
	DepositRecord bson.Raw `bson:"deposit_record"`
}

func (sv *DepositRecordStateValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of DepositRecordStateValue")

	var u DepositRecordStateValueBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	sv.BaseHinter = hint.NewBaseHinter(ht)

	var depositInfo types.DepositRecord
	if err := depositInfo.DecodeBSON(u.DepositRecord, enc); err != nil {
		return e.Wrap(err)
	}
	sv.Record = depositInfo

	return nil
}
