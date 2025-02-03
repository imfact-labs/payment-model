package payment

import (
	"context"
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var depositProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(DepositProcessor)
	},
}

func (Deposit) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type DepositProcessor struct {
	*base.BaseOperationProcessor
}

func NewDepositProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new DepositProcessor")

		nopp := depositProcessorPool.Get()
		opp, ok := nopp.(*DepositProcessor)
		if !ok {
			return nil, e.Errorf("expected DepositProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b

		return opp, nil
	}
}

func (opp *DepositProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(DepositFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", DepositFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	_, err := cstate.ExistsState(currency.BalanceStateKey(fact.Sender(), fact.Amount().Currency()),
		fmt.Sprintf("balance of account, %v", fact.Sender()), getStateFunc,
	)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if err := cstate.CheckExistsState(state.DesignStateKey(fact.Contract().String()), getStateFunc); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("payment service for contract account %v",
				fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *DepositProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(DepositFact)

	cid := fact.Amount().Currency()
	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, err := state.GetDesignFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError("service design value not found, %q; %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc
	accountInfo := design.Account(fact.Sender().String())
	if accountInfo != nil {
		// additional deposit
		st, _ := cstate.ExistsState(state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()), "account record", getStateFunc)
		accountRecord, _ := state.GetAccountRecordFromState(st)
		nAccountRecord := types.NewAccountRecord(fact.Sender())
		nAmount := ctypes.NewAmount(accountRecord.Amount(cid.String()).Big().Add(fact.Amount().Big()), cid)
		nAccountRecord.SetAmount(cid.String(), nAmount)
		nAccountRecord.SetLastTime(cid.String(), *accountRecord.LastTime(cid.String()))

		if err := nAccountRecord.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid record of account, %v in contract account, %v: %w", fact.Sender(), fact.Contract(), err), nil
		}
		// update AccountRecord
		sts = append(sts, cstate.NewStateMergeValue(
			state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()),
			state.NewAccountRecordStateValue(nAccountRecord),
		))
	} else {
		// new deposit
		accountInfo := types.NewAccountInfo(fact.Sender())
		amount := ctypes.NewAmount(fact.TransferLimit(), fact.Amount().Currency())
		accountInfo.SetTransferLimit(amount)
		accountInfo.SetPeriodTime(
			fact.Amount().Currency().String(),
			[3]uint64{fact.startTime, fact.endTime, fact.duration},
		)
		err = design.AddAccount(accountInfo)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to add info of account, %v in contract account %v: %w", fact.Sender(), fact.Contract(), err,
			), nil
		}
		if err := design.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
		}

		sts = append(sts, cstate.NewStateMergeValue(
			state.DesignStateKey(fact.Contract().String()),
			state.NewDesignStateValue(design),
		))

		// new AccountRecord
		nAccountRecord := types.NewAccountRecord(fact.Sender())
		nAccountRecord.SetAmount(cid.String(), fact.Amount())
		nAccountRecord.SetLastTime(cid.String(), 0)

		if err := nAccountRecord.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid record of account, %v in contract account, %v; %w", fact.Sender(), fact.Contract(), err), nil
		}

		sts = append(sts, cstate.NewStateMergeValue(
			state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()),
			state.NewAccountRecordStateValue(nAccountRecord),
		))

	}

	sts = append(
		sts,
		common.NewBaseStateMergeValue(
			currency.BalanceStateKey(fact.Sender(), fact.Amount().Currency()),
			currency.NewDeductBalanceStateValue(fact.Amount()),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(
					height, currency.BalanceStateKey(fact.Sender(), fact.Amount().Currency()),
					fact.Amount().Currency(), st,
				)
			}),
	)

	sts = append(sts, common.NewBaseStateMergeValue(
		currency.BalanceStateKey(fact.Contract(), fact.Amount().Currency()),
		currency.NewAddBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return currency.NewBalanceStateValueMerger(height,
				currency.BalanceStateKey(fact.Contract(), fact.Amount().Currency()),
				fact.Amount().Currency(), st,
			)
		},
	))

	return sts, nil, nil
}

func (opp *DepositProcessor) Close() error {
	depositProcessorPool.Put(opp)

	return nil
}
