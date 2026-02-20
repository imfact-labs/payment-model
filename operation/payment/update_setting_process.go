package payment

import (
	"context"
	"sync"

	"github.com/imfact-labs/currency-model/common"
	cstate "github.com/imfact-labs/currency-model/state"
	ctypes "github.com/imfact-labs/currency-model/types"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/payment-model/state"
	"github.com/imfact-labs/payment-model/types"
)

var updateAccountSettingProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(UpdateAccountSettingProcessor)
	},
}

func (UpdateAccountSetting) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type UpdateAccountSettingProcessor struct {
	*base.BaseOperationProcessor
}

func NewUpdateAccountSettingProcessor() ctypes.GetNewProcessor {
	return func(
		height base.Height,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new UpdateAccountSettingProcessor")

		nopp := updateAccountSettingProcessorPool.Get()
		opp, ok := nopp.(*UpdateAccountSettingProcessor)
		if !ok {
			return nil, e.Errorf("expected UpdateAccountSettingProcessor, not %T", nopp)
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

func (opp *UpdateAccountSettingProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(UpdateAccountSettingFact)

	cid := fact.Currency()
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", UpdateAccountSettingFact{}, op.Fact())), nil
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
				Wrap(common.ErrMServiceNF).Errorf("payment service state for contract account %v",
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

	setting := design.AccountSetting(fact.Sender().String())
	if setting == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf("setting of account, %v not found in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}

	big := setting.TransferLimit(cid.String())
	if big == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf("setting for currency, %v of account, %v not found in contract account %v",
				cid, fact.Sender(), fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *UpdateAccountSettingProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(UpdateAccountSettingFact)

	cid := fact.Currency()
	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, _ := state.GetDesignFromState(st)
	setting := design.AccountSetting(fact.Sender().String())
	nSetting := types.NewSettings(fact.Sender())
	for k, v := range setting.Items() {
		nSetting.SetItem(k, v.TransferLimit, v.StartTime, v.EndTime, v.Duration)
	}
	nSetting.SetItem(cid.String(), fact.TransferLimit(), fact.StartTime(), fact.EndTime(), fact.Duration())

	nDesign := types.NewDesign()
	for _, v := range design.AccountSettings() {
		nDesign.AddAccountSetting(v)
	}
	if err := nDesign.UpdateAccountSetting(nSetting); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"failed to update setting of account, %v in contract account, %v: %w", fact.Sender(), fact.Contract(), err), nil
	}

	if err := nDesign.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc
	sts = append(sts, cstate.NewStateMergeValue(
		state.DesignStateKey(fact.Contract().String()),
		state.NewDesignStateValue(nDesign),
	))

	return sts, nil, nil
}

func (opp *UpdateAccountSettingProcessor) Close() error {
	updateAccountSettingProcessorPool.Put(opp)

	return nil
}
