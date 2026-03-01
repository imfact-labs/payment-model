package spec

import (
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/payment-model/operation/payment"
	"github.com/imfact-labs/payment-model/state"
	"github.com/imfact-labs/payment-model/types"
)

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit

	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.SettingHint, Instance: types.Setting{}},
	{Hint: types.DepositRecordHint, Instance: types.DepositRecord{}},

	{Hint: payment.DepositHint, Instance: payment.Deposit{}},
	{Hint: payment.RegisterModelHint, Instance: payment.RegisterModel{}},
	{Hint: payment.TransferHint, Instance: payment.Transfer{}},
	{Hint: payment.UpdateAccountSettingHint, Instance: payment.UpdateAccountSetting{}},
	{Hint: payment.WithdrawHint, Instance: payment.Withdraw{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.DepositRecordStateValueHint, Instance: state.DepositRecordStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: payment.DepositFactHint, Instance: payment.DepositFact{}},
	{Hint: payment.RegisterModelFactHint, Instance: payment.RegisterModelFact{}},
	{Hint: payment.TransferFactHint, Instance: payment.TransferFact{}},
	{Hint: payment.UpdateAccountSettingFactHint, Instance: payment.UpdateAccountSettingFact{}},
	{Hint: payment.WithdrawFactHint, Instance: payment.WithdrawFact{}},
}
