package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var PaymentIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "payment_contract_height"),
	},
}

var PaymentAccountRecordIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "address", Value: 1},
			bson.E{Key: "height", Value: -1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "payment_contract_address_height"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNamePayment] = PaymentIndexModels
	DefaultIndexes[DefaultColNamePaymentAccount] = PaymentAccountRecordIndexModels
}
