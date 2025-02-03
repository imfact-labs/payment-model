package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

type DesignDoc struct {
	mongodb.BaseDoc
	st     base.State
	design types.Design
}

// NewDesignDoc get the State of TimeStamp Design
func NewDesignDoc(st base.State, enc encoder.Encoder) (DesignDoc, error) {
	design, err := state.GetDesignFromState(st)

	if err != nil {
		return DesignDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DesignDoc{}, err
	}

	return DesignDoc{
		BaseDoc: b,
		st:      st,
		design:  design,
	}, nil
}

func (doc DesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.PaymentStateKeyPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type AccountRecordDoc struct {
	mongodb.BaseDoc
	st            base.State
	accountRecord types.AccountRecord
}

func NewAccountRecordDoc(st base.State, enc encoder.Encoder) (*AccountRecordDoc, error) {
	accountRecord, err := state.GetAccountRecordFromState(st)
	if err != nil {
		return nil, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &AccountRecordDoc{
		BaseDoc:       b,
		st:            st,
		accountRecord: *accountRecord,
	}, nil
}

func (doc AccountRecordDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.PaymentStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["address"] = doc.accountRecord.Address()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

var (
	AccountInfoValueHint = hint.MustNewHint("mitum-payment-account-info-value-v0.0.1")
)

type AccountInfoValue struct {
	hint.BaseHinter
	accountInfo   types.AccountInfo
	accountRecord types.AccountRecord
}

func NewAccountInfoValue(
	accountInfo types.AccountInfo,
	accountRecord types.AccountRecord,
) AccountInfoValue {
	return AccountInfoValue{
		BaseHinter:    hint.NewBaseHinter(AccountInfoValueHint),
		accountInfo:   accountInfo,
		accountRecord: accountRecord,
	}
}

func (ai AccountInfoValue) Hint() hint.Hint {
	return AccountInfoValueHint
}

func (ai AccountInfoValue) AccountInfo() types.AccountInfo {
	return ai.accountInfo
}

func (ai AccountInfoValue) AccountRecord() types.AccountRecord {
	return ai.accountRecord
}

//func (ai AccountInfoValue) MarshalBSON() ([]byte, error) {
//	accountInfo, err := ai.accountInfo.MarshalBSON()
//	if err != nil {
//		return nil, err
//	}
//
//	accountRecord, err := ai.accountRecord.MarshalBSON()
//	if err != nil {
//		return nil, err
//	}
//
//	return bsonenc.Marshal(
//		bson.M{
//			"_hint":          ai.Hint().String(),
//			"account_info":   accountInfo,
//			"account_record": accountRecord,
//		},
//	)
//}
//
//type AccountInfoValueBSONUnmarshaler struct {
//	Hint          string   `bson:"_hint"`
//	AccountInfo   bson.Raw `bson:"account_info"`
//	AccountRecord bson.Raw `bson:"account_record"`
//}
//
//func (ai *AccountInfoValue) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
//	e := util.StringError("Decode bson of AccountInfoValue")
//
//	var uai AccountInfoValueBSONUnmarshaler
//	if err := enc.Unmarshal(b, &uai); err != nil {
//		return e.Wrap(err)
//	}
//
//	ht, err := hint.ParseHint(uai.Hint)
//	if err != nil {
//		return e.Wrap(err)
//	}
//
//	ai.BaseHinter = hint.NewBaseHinter(ht)
//
//	var accountInfo types.AccountInfo
//	if err := accountInfo.DecodeBSON(uai.AccountInfo, enc); err != nil {
//		return e.Wrap(err)
//	}
//
//	var accountRecord types.AccountRecord
//	if err := accountRecord.DecodeBSON(uai.AccountRecord, enc); err != nil {
//		return e.Wrap(err)
//	}
//
//	ai.accountInfo = accountInfo
//	ai.accountRecord = accountRecord
//
//	return nil
//}

type AccountInfoValueJSONMarshaler struct {
	hint.BaseHinter
	AccountInfo   types.AccountInfo   `json:"account_info"`
	AccountRecord types.AccountRecord `json:"account_record"`
}

func (ai AccountInfoValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AccountInfoValueJSONMarshaler{
		BaseHinter:    ai.BaseHinter,
		AccountInfo:   ai.accountInfo,
		AccountRecord: ai.accountRecord,
	})
}
