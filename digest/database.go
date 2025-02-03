package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	utilc "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultColNamePayment        = "digest_pmt"
	DefaultColNamePaymentAccount = "digest_pmt_ac"
)

func PaymentDesign(db *cdigest.Database, contract string) (*types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := db.MongoClient().GetByFilter(
		DefaultColNamePayment,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return nil, nil, utilm.ErrNotFound.WithMessage(err, "payment design by contract account %v", contract)
	}

	if sta != nil {
		de, err := state.GetDesignFromState(sta)
		if err != nil {
			return nil, nil, err
		}
		return &de, sta, nil
	} else {
		return nil, nil, errors.Errorf("state is nil")
	}
}

func AccountInfo(db *cdigest.Database, contract, account string) (*AccountInfoValue, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("address", account)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var st base.State
	var design *types.Design
	var accountRecord *types.AccountRecord
	var err error
	if err := db.MongoClient().GetByFilter(
		DefaultColNamePaymentAccount,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			st = i
			return nil
		},
		opt,
	); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			emptyAccountRecord := types.NewEmptyAccountRecord()
			accountRecord = &emptyAccountRecord
		} else {
			return nil,
				utilm.ErrNotFound.WithMessage(
					err, "payment account record by contract account %s, account %s", contract, account,
				)
		}
	} else {
		if st != nil {
			accountRecord, err = state.GetAccountRecordFromState(st)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.Errorf("state is nil")
		}
	}
	if design, _, err = PaymentDesign(db, contract); err != nil {
		return nil, err
	}

	accountInfo := design.Account(account)
	if accountInfo == nil {
		return nil,
			utilm.ErrNotFound.WithMessage(
				err, "payment account info by contract account %s, account %s", contract, account,
			)
	}

	accountInfoValue := NewAccountInfoValue(*accountInfo, *accountRecord)
	return &accountInfoValue, nil
}
