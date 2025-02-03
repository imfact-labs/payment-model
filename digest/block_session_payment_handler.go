package digest

import (
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) preparePayment() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var paymentDesignModels []mongo.WriteModel
	var paymentAccountModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsDesignStateKey(st.Key()):
			j, err := bs.handlePaymentDesignState(st)
			if err != nil {
				return err
			}
			paymentDesignModels = append(paymentDesignModels, j...)
		case state.IsAccountRecordStateKey(st.Key()):
			j, err := bs.handlePaymentAccountRecordState(st)
			if err != nil {
				return err
			}
			paymentAccountModels = append(paymentAccountModels, j...)
		default:
			continue
		}
	}

	bs.paymentDesignModels = paymentDesignModels
	bs.paymentAccountModels = paymentAccountModels

	return nil
}

func (bs *BlockSession) handlePaymentDesignState(st base.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewDesignDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handlePaymentAccountRecordState(st base.State) ([]mongo.WriteModel, error) {
	if AccountRecordDoc, err := NewAccountRecordDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(AccountRecordDoc),
		}, nil
	}
}
