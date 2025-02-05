package cmds

import (
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-payment/operation/deposit"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/pkg/errors"
)

var Hinters []encoder.DecodeDetail
var SupportedProposalOperationFactHinters []encoder.DecodeDetail

var AddedHinters = []encoder.DecodeDetail{
	// revive:disable-next-line:line-length-limit

	{Hint: types.DesignHint, Instance: types.Design{}},
	{Hint: types.SettingHint, Instance: types.Setting{}},
	{Hint: types.DepositRecordHint, Instance: types.DepositRecord{}},

	{Hint: deposit.DepositHint, Instance: deposit.Deposit{}},
	{Hint: deposit.RegisterModelHint, Instance: deposit.RegisterModel{}},
	{Hint: deposit.TransferHint, Instance: deposit.Transfer{}},
	{Hint: deposit.WithdrawHint, Instance: deposit.Withdraw{}},

	{Hint: state.DesignStateValueHint, Instance: state.DesignStateValue{}},
	{Hint: state.DepositRecordStateValueHint, Instance: state.DepositRecordStateValue{}},
}

var AddedSupportedHinters = []encoder.DecodeDetail{
	{Hint: deposit.DepositFactHint, Instance: deposit.DepositFact{}},
	{Hint: deposit.RegisterModelFactHint, Instance: deposit.RegisterModelFact{}},
	{Hint: deposit.TransferFactHint, Instance: deposit.TransferFact{}},
	{Hint: deposit.WithdrawFactHint, Instance: deposit.WithdrawFact{}},
}

func init() {
	Hinters = append(Hinters, currencycmds.Hinters...)
	Hinters = append(Hinters, AddedHinters...)

	SupportedProposalOperationFactHinters = append(SupportedProposalOperationFactHinters, currencycmds.SupportedProposalOperationFactHinters...)
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
