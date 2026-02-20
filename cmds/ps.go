package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	cprocessor "github.com/imfact-labs/currency-model/operation/processor"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/isaac"
	"github.com/imfact-labs/mitum2/launch"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/mitum2/util/hint"
	"github.com/imfact-labs/mitum2/util/ps"
	"github.com/imfact-labs/payment-model/operation/payment"
)

var PNameOperationProcessorsMap = ps.Name("mitum-payment-operation-processors-map")

func POperationProcessorsMap(pctx context.Context) (context.Context, error) {
	var isaacParams *isaac.Params
	var db isaac.Database
	var opr *cprocessor.OperationProcessor
	var setA *hint.CompatibleSet[isaac.NewOperationProcessorInternalFunc]
	var setB *hint.CompatibleSet[ccmds.NewOperationProcessorInternalWithProposalFunc]

	if err := util.LoadFromContextOK(pctx,
		launch.ISAACParamsContextKey, &isaacParams,
		launch.CenterDatabaseContextKey, &db,
		ccmds.OperationProcessorContextKey, &opr,
		launch.OperationProcessorsMapContextKey, &setA,
		ccmds.OperationProcessorsMapBContextKey, &setB,
	); err != nil {
		return pctx, err
	}

	//err := opr.SetCheckDuplicationFunc(processor.CheckDuplication)
	//if err != nil {
	//	return pctx, err
	//}
	err := opr.SetGetNewProcessorFunc(cprocessor.GetNewProcessor)
	if err != nil {
		return pctx, err
	}

	if err := opr.SetProcessor(
		payment.RegisterModelHint,
		payment.NewRegisterModelProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		payment.DepositHint,
		payment.NewDepositProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessorWithProposal(
		payment.WithdrawHint,
		payment.NewWithdrawProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessor(
		payment.UpdateAccountSettingHint,
		payment.NewUpdateAccountSettingProcessor(),
	); err != nil {
		return pctx, err
	} else if err := opr.SetProcessorWithProposal(
		payment.TransferHint,
		payment.NewTransferProcessor(),
	); err != nil {
		return pctx, err
	}

	_ = setA.Add(payment.RegisterModelHint, func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	_ = setA.Add(payment.DepositHint, func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	_ = setB.Add(payment.WithdrawHint, func(height base.Height, proposal base.ProposalSignFact, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		if err := opr.SetProposal(&proposal); err != nil {
			return nil, err
		}
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	_ = setA.Add(payment.UpdateAccountSettingHint, func(height base.Height, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	_ = setB.Add(payment.TransferHint, func(height base.Height, proposal base.ProposalSignFact, getStatef base.GetStateFunc) (base.OperationProcessor, error) {
		if err := opr.SetProposal(&proposal); err != nil {
			return nil, err
		}
		return opr.New(
			height,
			getStatef,
			nil,
			nil,
		)
	})

	pctx = context.WithValue(pctx, ccmds.OperationProcessorContextKey, opr)
	pctx = context.WithValue(pctx, launch.OperationProcessorsMapContextKey, setA) //revive:disable-line:modifies-parameter
	pctx = context.WithValue(pctx, ccmds.OperationProcessorsMapBContextKey, setB) //revive:disable-line:modifies-parameter

	return pctx, nil
}
