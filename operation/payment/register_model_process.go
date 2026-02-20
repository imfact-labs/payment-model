package payment

import (
	"context"
	"sync"

	"github.com/imfact-labs/currency-model/common"
	cstate "github.com/imfact-labs/currency-model/state"
	cestate "github.com/imfact-labs/currency-model/state/extension"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/payment-model/state"
	"github.com/imfact-labs/payment-model/types"
	"github.com/pkg/errors"
)

var registerModelProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(RegisterModelProcessor)
	},
}

type RegisterModelProcessor struct {
	*base.BaseOperationProcessor
}

func NewRegisterModelProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new CreateServiceProcessor")

		nopp := registerModelProcessorPool.Get()
		opp, ok := nopp.(*RegisterModelProcessor)
		if !ok {
			return nil, errors.Errorf("expected RegisterModelProcessor, not %T", nopp)
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

func (opp *RegisterModelProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(RegisterModelFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", RegisterModelFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	if found, _ := cstate.CheckNotExistsState(state.DesignStateKey(fact.Contract().String()), getStateFunc); found {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceE).Errorf("timestamp service for contract account %v",
				fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *RegisterModelProcessor) Process(
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(RegisterModelFact)

	var sts []base.StateMergeValue

	design := types.NewDesign()
	if err := design.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid timestamp design, %q; %w", fact.Contract(), err), nil
	}

	sts = append(sts, cstate.NewStateMergeValue(
		state.DesignStateKey(fact.Contract().String()),
		state.NewDesignStateValue(design),
	))

	st, _ := cstate.ExistsState(cestate.StateKeyContractAccount(fact.Contract()), "contract account", getStateFunc)
	ca, _ := cestate.StateContractAccountValue(st)
	ca.SetActive(true)
	h := op.Hint()
	ca.SetRegisterOperation(&h)
	ca.SetBalanceStatus(ctypes.WithdrawalBlocked)

	sts = append(sts, cstate.NewStateMergeValue(
		cestate.StateKeyContractAccount(fact.Contract()),
		cestate.NewContractAccountStateValue(ca),
	))

	return sts, nil, nil
}

func (opp *RegisterModelProcessor) Close() error {
	registerModelProcessorPool.Put(opp)

	return nil
}
