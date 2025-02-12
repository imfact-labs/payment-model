package payment

import (
	"context"
	"github.com/ProtoconNet/mitum-currency/v3/state/currency"
	"github.com/ProtoconNet/mitum-payment/types"
	"sync"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	ctypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-payment/state"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
)

var withdrawProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(WithdrawProcessor)
	},
}

func (Withdraw) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type WithdrawProcessor struct {
	*base.BaseOperationProcessor
	proposal *base.ProposalSignFact
}

func NewWithdrawProcessor() ctypes.GetNewProcessorWithProposal {
	return func(
		height base.Height,
		proposal *base.ProposalSignFact,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new WithdrawProcessor")

		nopp := withdrawProcessorPool.Get()
		opp, ok := nopp.(*WithdrawProcessor)
		if !ok {
			return nil, e.Errorf("expected WithdrawProcessor, not %T", nopp)
		}

		b, err := base.NewBaseOperationProcessor(
			height, getStateFunc, newPreProcessConstraintFunc, newProcessConstraintFunc)
		if err != nil {
			return nil, e.Wrap(err)
		}

		opp.BaseOperationProcessor = b
		opp.proposal = proposal

		return opp, nil
	}
}

func (opp *WithdrawProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(WithdrawFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", WithdrawFact{}, op.Fact())), nil
	}

	cid := fact.Currency()
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

	st, err = cstate.ExistsState(
		state.DepositRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		"account record", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateNF).Errorf("record of account, %v in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}

	record, err := state.GetDepositRecordFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateValInvalid).Errorf("record of account, %v not found in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}
	amount := record.Amount(cid.String())
	if amount == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf(
				"record of account, %v for currency id, %v not found in contract account %v",
				fact.Sender(), cid, fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *WithdrawProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(WithdrawFact)

	var sts []base.StateMergeValue // nolint:prealloc
	proposal := *opp.proposal
	nowTime := uint64(proposal.ProposalFact().ProposedAt().Unix())
	cid := fact.Currency()
	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, _ := state.GetDesignFromState(st)
	setting := design.AccountSetting(fact.Sender().String())
	nSetting := types.NewSettings(fact.Sender())
	for k, v := range setting.Items() {
		nSetting.SetItem(k, v.TransferLimit, v.StartTime, v.EndTime, v.Duration)
	}
	nSetting.Remove(cid.String())
	nDesign := types.NewDesign()
	for _, v := range design.AccountSettings() {
		nDesign.AddAccountSetting(v)
	}
	if len(nSetting.Items()) < 1 {
		nDesign.RemoveAccountSetting(fact.Sender())
		if err := nDesign.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
		}
	} else {
		nDesign.UpdateAccountSetting(nSetting)
		if err := nDesign.IsValid(nil); err != nil {
			return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
		}
	}
	// update Design
	sts = append(sts, cstate.NewStateMergeValue(
		state.DesignStateKey(fact.Contract().String()),
		state.NewDesignStateValue(nDesign),
	))

	st, _ = cstate.ExistsState(
		state.DepositRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		"account record", getStateFunc)
	record, _ := state.GetDepositRecordFromState(st)
	big := record.Amount(cid.String())
	nRecord := types.NewDepositRecord(fact.Sender())
	for k, v := range record.Items() {
		nRecord.SetItem(k, v.Amount, v.TransferredAt)
	}
	nRecord.SetItem(cid.String(), common.ZeroBig, nowTime)

	if err := nRecord.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			"invalid record of account, %v in contract account, %v: %w", fact.Sender(), fact.Contract(), err), nil
	}
	// update Record
	sts = append(sts, cstate.NewStateMergeValue(
		state.DepositRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		state.NewDepositRecordStateValue(nRecord),
	))

	am := ctypes.NewAmount(*big, cid)
	sts = append(
		sts,
		common.NewBaseStateMergeValue(
			currency.BalanceStateKey(fact.Contract(), cid),
			currency.NewDeductBalanceStateValue(am),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(
					height, currency.BalanceStateKey(fact.Contract(), cid),
					cid, st,
				)
			}),
	)

	sts = append(sts, common.NewBaseStateMergeValue(
		currency.BalanceStateKey(fact.Sender(), cid),
		currency.NewAddBalanceStateValue(am),
		func(height base.Height, st base.State) base.StateValueMerger {
			return currency.NewBalanceStateValueMerger(height,
				currency.BalanceStateKey(fact.Sender(), cid),
				cid, st,
			)
		},
	))

	return sts, nil, nil
}

func (opp *WithdrawProcessor) Close() error {
	opp.proposal = nil
	withdrawProcessorPool.Put(opp)

	return nil
}
