package digest

import (
	currencydigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func PreparePayment(bs *currencydigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
	switch {
	case state.IsDesignStateKey(st.Key()):
		j, err := handlePaymentDesignState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNamePayment, j, nil
	case state.IsDepositRecordStateKey(st.Key()):
		j, err := handlePaymentAccountRecordState(bs, st)
		if err != nil {
			return "", nil, err
		}

		return DefaultColNamePaymentAccount, j, nil
	}

	return "", nil, nil
}

func handlePaymentDesignState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewDesignDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func handlePaymentAccountRecordState(bs *currencydigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if AccountRecordDoc, err := NewDepositRecordDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(AccountRecordDoc),
		}, nil
	}
}
