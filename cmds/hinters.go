package cmds

import (
	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/util/encoder"
	"github.com/imfact-labs/payment-model/operation/payment"
	"github.com/imfact-labs/payment-model/state"
	"github.com/imfact-labs/payment-model/types"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

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

func init() {
	Hinters = append(Hinters, ccmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, ccmds.SupportedProposalOperationFactHinters...)
	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, AddedSupportedHinters...)
}

func LoadHinters(encs *encoder.Encoders) error {
	for i := range Hinters {
		if err := encs.AddDetail(Hinters[i]); err != nil {
			return errors.Wrap(err, "add hinter to encoder")
		}
	}

	for i := range SupportedProposalOperationFactHinters {
		if err := encs.AddDetail(SupportedProposalOperationFactHinters[i]); err != nil {
			return errors.Wrap(err, "add supported proposal operation fact hinter to encoder")
		}
	}

	return nil
}
