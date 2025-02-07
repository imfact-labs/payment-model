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

	_, err := cstate.ExistsState(currency.BalanceStateKey(fact.Sender(), fact.Currency()),
		fmt.Sprintf("balance of account, %v", fact.Sender()), getStateFunc,
	)
	if err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	st, err := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMServiceNF).Errorf("payment service for contract account %v",
				fact.Contract(),
			)), nil
	}
	design, err := state.GetDesignFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateValInvalid).Errorf(
				"service design value not found, %v: %v", fact.Contract(), err)), nil
	}
	setting := design.AccountSetting(fact.Sender().String())
	if setting != nil {
		st, err = cstate.ExistsState(state.DepositRecordStateKey(
			fact.Contract().String(), fact.Sender().String()), "account record", getStateFunc)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				common.ErrMPreProcess.
					Wrap(common.ErrMStateNF).Errorf(
					"record of account, %v nof found in contract account, %v: %v", fact.Sender(), fact.Contract(), err)), nil
		}
	}

	return ctx, nil, nil
}

func (opp *DepositProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(DepositFact)

	cid := fact.Currency()
	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, _ := state.GetDesignFromState(st)

	var sts []base.StateMergeValue // nolint:prealloc
	setting := design.AccountSetting(fact.Sender().String())

	if setting != nil {
		// additional deposit
		st, _ := cstate.ExistsState(state.DepositRecordStateKey(
			fact.Contract().String(), fact.Sender().String()), "account record", getStateFunc)
		record, _ := state.GetDepositRecordFromState(st)
		amount := record.Amount(cid.String())

		var nAmount common.Big
		var nTransfferdAt uint64
		if amount == nil {
			nAmount = fact.Amount()
			nTransfferdAt = 0
		} else {
			nAmount = amount.Add(fact.Amount())
			nTransfferdAt = *record.TransferredAt(cid.String())
		}

		nRecord := types.NewDepositRecord(fact.Sender())
		for k, v := range record.Items() {
			nRecord.SetItem(k, v.Amount, v.TransferredAt)
		}
		nRecord.SetItem(cid.String(), nAmount, nTransfferdAt)

		if err := nRecord.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"invalid record of account, %v in contract account, %v: %w", fact.Sender(), fact.Contract(), err), nil
		}
		// update Record
		sts = append(sts, cstate.NewStateMergeValue(
			state.DepositRecordStateKey(fact.Contract().String(), fact.Sender().String()),
			state.NewDepositRecordStateValue(nRecord),
		))

		// update AccountSetting
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

		sts = append(sts, cstate.NewStateMergeValue(
			state.DesignStateKey(fact.Contract().String()),
			state.NewDesignStateValue(nDesign),
		))
	} else {
		nSetting := types.NewSettings(fact.Sender())
		nSetting.SetItem(cid.String(), fact.TransferLimit(), fact.StartTime(), fact.EndTime(), fact.Duration())
		nDesign := types.NewDesign()
		for _, v := range design.AccountSettings() {
			nDesign.AddAccountSetting(v)
		}
		err := nDesign.AddAccountSetting(nSetting)
		if err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"failed to add setting of account, %v in contract account %v: %w", fact.Sender(), fact.Contract(), err,
			), nil
		}
		if err := nDesign.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
		}

		sts = append(sts, cstate.NewStateMergeValue(
			state.DesignStateKey(fact.Contract().String()),
			state.NewDesignStateValue(nDesign),
		))

		// new AccountRecord
		nRecord := types.NewDepositRecord(fact.Sender())
		nRecord.SetItem(cid.String(), fact.Amount(), 0)

		if err := nRecord.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError(
				"invalid record of account, %v in contract account, %v; %w", fact.Sender(), fact.Contract(), err), nil
		}

		sts = append(sts, cstate.NewStateMergeValue(
			state.DepositRecordStateKey(fact.Contract().String(), fact.Sender().String()),
			state.NewDepositRecordStateValue(nRecord),
		))
	}

	am := ctypes.NewAmount(fact.Amount(), cid)
	sts = append(
		sts,
		common.NewBaseStateMergeValue(
			currency.BalanceStateKey(fact.Sender(), cid),
			currency.NewDeductBalanceStateValue(am),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(
					height, currency.BalanceStateKey(fact.Sender(), cid),
					cid, st,
				)
			}),
	)

	sts = append(sts, common.NewBaseStateMergeValue(
		currency.BalanceStateKey(fact.Contract(), cid),
		currency.NewAddBalanceStateValue(am),
		func(height base.Height, st base.State) base.StateValueMerger {
			return currency.NewBalanceStateValueMerger(height,
				currency.BalanceStateKey(fact.Contract(), cid),
				cid, st,
			)
		},
	))

	return sts, nil, nil
}

func (opp *DepositProcessor) Close() error {
	depositProcessorPool.Put(opp)

	return nil
}
