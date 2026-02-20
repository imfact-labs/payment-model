package digest

import (
	mongodb "github.com/imfact-labs/currency-model/digest/mongodb"
	cstate "github.com/imfact-labs/currency-model/state"
	"github.com/imfact-labs/currency-model/utils/bsonenc"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/payment-model/state"
	"github.com/imfact-labs/payment-model/types"
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

type DepositRecordDoc struct {
	mongodb.BaseDoc
	st     base.State
	record types.DepositRecord
}

func NewDepositRecordDoc(st base.State, enc encoder.Encoder) (*DepositRecordDoc, error) {
	record, err := state.GetDepositRecordFromState(st)
	if err != nil {
		return nil, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return nil, err
	}

	return &DepositRecordDoc{
		BaseDoc: b,
		st:      st,
		record:  *record,
	}, nil
}

func (doc DepositRecordDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.PaymentStateKeyPrefix, 4)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["address"] = doc.record.Address()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

var (
	AccountInfoValueHint = hint.MustNewHint("mitum-payment-account-info-value-v0.0.1")
)

type AccountInfoValue struct {
	hint.BaseHinter
	setting types.Setting
	record  types.DepositRecord
}

func NewAccountInfoValue(
	setting types.Setting,
	record types.DepositRecord,
) AccountInfoValue {
	return AccountInfoValue{
		BaseHinter: hint.NewBaseHinter(AccountInfoValueHint),
		setting:    setting,
		record:     record,
	}
}

func (ai AccountInfoValue) Hint() hint.Hint {
	return AccountInfoValueHint
}

func (ai AccountInfoValue) AccountInfo() types.Setting {
	return ai.setting
}

func (ai AccountInfoValue) AccountRecord() types.DepositRecord {
	return ai.record
}

type AccountInfoValueJSONMarshaler struct {
	hint.BaseHinter
	Setting types.Setting       `json:"transfer_setting"`
	Record  types.DepositRecord `json:"deposit_record"`
}

func (ai AccountInfoValue) MarshalJSON() ([]byte, error) {
	return util.MarshalJSON(AccountInfoValueJSONMarshaler{
		BaseHinter: ai.BaseHinter,
		Setting:    ai.setting,
		Record:     ai.record,
	})
}
