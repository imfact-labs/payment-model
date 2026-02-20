package digest

import (
	cdigest "github.com/imfact-labs/currency-model/digest"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/payment-model/state"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func PreparePayment(bs *cdigest.BlockSession, st base.State) (string, []mongo.WriteModel, error) {
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

func handlePaymentDesignState(bs *cdigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewDesignDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func handlePaymentAccountRecordState(bs *cdigest.BlockSession, st base.State) ([]mongo.WriteModel, error) {
	if AccountRecordDoc, err := NewDepositRecordDoc(st, bs.Database().Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(AccountRecordDoc),
		}, nil
	}
}
