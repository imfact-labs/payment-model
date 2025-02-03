package payment

import (
	"context"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum-payment/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var updateAccountInfoProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateAccountInfoProcessor)
	},
}

func (UpdateAccountInfo) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateAccountInfoProcessor struct {
	*base.BaseOperationProcessor
}

func NewUpdateAccountInfoProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new UpdateAccountInfoProcessor")

		nopp := updateAccountInfoProcessorPool.Get()
		opp, ok := nopp.(*UpdateAccountInfoProcessor)
		if !ok {
			return nil, e.Errorf("expected UpdateAccountInfoProcessor, not %T", nopp)
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

func (opp *UpdateAccountInfoProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(UpdateAccountInfoFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", UpdateAccountInfoFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	st, err := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("service for contract account %v",
				fact.Contract(),
			)), nil
	}

	design, err := state.GetDesignFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateNF).Errorf("service design value for contract account %v",
				fact.Contract(),
			)), nil
	}

	accountInfo := design.Account(fact.Sender().String())
	if accountInfo == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf("info of account, %v not found in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateAccountInfoProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(UpdateAccountInfoFact)

	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, _ := state.GetDesignFromState(st)

	accountInfo := types.NewAccountInfo(fact.Sender())
	accountInfo.SetTransferLimit(fact.transferLimit)
	accountInfo.SetPeriodTime(
		fact.TransferLimit().Currency().String(),
		[3]uint64{fact.StartTime(), fact.EndTime(), fact.Duration()},
	)
	err := design.UpdateAccount(accountInfo)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to update info of account, %v in contract account %v: %w", fact.Sender(), fact.Contract(), err,
		), nil
	}
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc
	sts = append(sts, cstate.NewStateMergeValue(
		state.DesignStateKey(fact.Contract().String()),
		state.NewDesignStateValue(design),
	))

	return sts, nil, nil
}

func (opp *UpdateAccountInfoProcessor) Close() error {
	updateAccountInfoProcessorPool.Put(opp)

	return nil
}
