package deposit

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum-currency/v3/operation/extras"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
	"go.mongodb.org/mongo-driver/bson"
)

func (fact DepositFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":          fact.Hint().String(),
			"hash":           fact.BaseFact.Hash().String(),
			"token":          fact.BaseFact.Token(),
			"sender":         fact.sender,
			"contract":       fact.contract,
			"amount":         fact.amount,
			"transfer_limit": fact.transferLimit,
			"start_time":     fact.startTime,
			"end_time":       fact.endTime,
			"duration":       fact.duration,
			"currency":       fact.currency,
		},
	)
}

type DepositFactBSONUnmarshaler struct {
	Hint          string     `bson:"_hint"`
	Sender        string     `bson:"sender"`
	Contract      string     `bson:"contract"`
	Amount        common.Big `bson:"amount"`
	TransferLimit common.Big `bson:"transfer_limit"`
	StartTime     uint64     `bson:"start_time"`
	EndTime       uint64     `bson:"end_time"`
	Duration      uint64     `bson:"duration"`
	Currency      string     `bson:"currency"`
}

func (fact *DepositFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var u common.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf DepositFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)
	fact.amount = uf.Amount
	fact.transferLimit = uf.TransferLimit

	if err := fact.unpack(
		enc, uf.Sender, uf.Contract, uf.StartTime, uf.EndTime, uf.Duration, uf.Currency,
	); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *fact)
	}

	return nil
}

func (op Deposit) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *Deposit) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperation = ubo

	var ueo extras.BaseOperationExtensions
	if err := ueo.DecodeBSON(b, enc); err != nil {
		return common.DecorateError(err, common.ErrDecodeBson, *op)
	}

	op.BaseOperationExtensions = &ueo

	return nil
}
