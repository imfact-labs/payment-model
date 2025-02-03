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

var transferProcessorPool = sync.Pool{
	New: func() interface{} {
		return new(TransferProcessor)
	},
}

func (Transfer) Process(
	_ context.Context, _ base.GetStateFunc,
) ([]base.StateMergeValue, base.OperationProcessReasonError, error) {
	return nil, nil, nil
}

type TransferProcessor struct {
	*base.BaseOperationProcessor
	proposal *base.ProposalSignFact
}

func NewTransferProcessor() ctypes.GetNewProcessorWithProposal {
	return func(
		height base.Height,
		proposal *base.ProposalSignFact,
		getStateFunc base.GetStateFunc,
		newPreProcessConstraintFunc base.NewOperationProcessorProcessFunc,
		newProcessConstraintFunc base.NewOperationProcessorProcessFunc,
	) (base.OperationProcessor, error) {
		e := util.StringError("failed to create new TransferProcessor")

		nopp := transferProcessorPool.Get()
		opp, ok := nopp.(*TransferProcessor)
		if !ok {
			return nil, e.Errorf("expected TransferProcessor, not %T", nopp)
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

func (opp *TransferProcessor) PreProcess(
	ctx context.Context, op base.Operation, getStateFunc base.GetStateFunc,
) (context.Context, base.OperationProcessReasonError, error) {
	fact, ok := op.Fact().(TransferFact)
	if !ok {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMTypeMismatch).
				Errorf("expected %T, not %T", TransferFact{}, op.Fact())), nil
	}

	if err := fact.IsValid(nil); err != nil {
		return ctx, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Errorf("%v", err)), nil
	}

	cid := fact.Amount().Currency()
	_, err := cstate.ExistsState(currency.BalanceStateKey(fact.Contract(), cid),
		fmt.Sprintf("balance of account, %v", fact.Contract()), getStateFunc,
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
	} else if tLimit := accountInfo.TransferLimit(cid.String()); tLimit.Big().Compare(fact.Amount().Big()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf(
				"transfer amount(%v) exceeds the limit(%v) of account, %v in contract account %v.",
				fact.Amount().Big(), tLimit.Big(), fact.Sender(), fact.Contract(),
			)), nil
	}

	st, err = cstate.ExistsState(
		state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		"account record", getStateFunc)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateNF).Errorf("record of account, %v in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}
	accountRecord, err := state.GetAccountRecordFromState(st)
	if err != nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMStateValInvalid).Errorf("record of account, %v in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	}
	amount := accountRecord.Amount(fact.Amount().Currency().String())
	if amount == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf("deposit of account, %v in contract account %v",
				fact.Sender(), fact.Contract(),
			)), nil
	} else if amount.Big().Compare(fact.Amount().Big()) < 0 {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf("transfer amount(%v) exceeds the deposit(%v) of account %v in contract account %v",
				fact.Amount().Big(), amount.Big(), fact.Sender(), fact.Contract(),
			)), nil
	} else if lastTime := accountRecord.LastTime(fact.Amount().Currency().String()); lastTime == nil {
		return nil, base.NewBaseOperationProcessReasonError(
			common.ErrMPreProcess.
				Wrap(common.ErrMValueInvalid).Errorf(
				"last time of transfer, %v in contract account %v.",
				fact.Sender(), fact.Contract(),
			)), nil
	}

	return ctx, nil, nil
}

func (opp *TransferProcessor) Process( // nolint:dupl
	_ context.Context, op base.Operation, getStateFunc base.GetStateFunc) (
	[]base.StateMergeValue, base.OperationProcessReasonError, error,
) {
	fact, _ := op.Fact().(TransferFact)

	cid := fact.Amount().Currency()
	proposal := *opp.proposal
	nowTime := uint64(proposal.ProposalFact().ProposedAt().Unix())
	var pTime *[3]uint64

	st, _ := cstate.ExistsState(state.DesignStateKey(fact.Contract().String()), "service design", getStateFunc)
	design, _ := state.GetDesignFromState(st)
	accountInfo := design.Account(fact.Sender().String())
	pTime = accountInfo.PeriodTime(fact.Amount().Currency().String())

	if pTime[0] > nowTime {
		return nil, base.NewBaseOperationProcessReasonError(
			"current time, %v is earlier than start time, %v for account, %v in contract account %v.",
			nowTime, pTime[0], fact.Sender(), fact.Contract(),
		), nil
	} else if pTime[1] < nowTime {
		return nil, base.NewBaseOperationProcessReasonError(
			"current time, %v is beyond the end time, %v for account, %v in contract account %v.",
			nowTime, pTime[1], fact.Sender(), fact.Contract(),
		), nil
	}

	st, _ = cstate.ExistsState(
		state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		"account record", getStateFunc)
	accountRecord, _ := state.GetAccountRecordFromState(st)
	if lastTime := accountRecord.LastTime(fact.Amount().Currency().String()); (*lastTime + pTime[2]) > nowTime {
		return nil, base.NewBaseOperationProcessReasonError(
			"last time of transfer, %v is too recent. Wait until required cool time, %v for account, %v in contract account %v.",
			*lastTime, pTime[2], fact.Sender(), fact.Contract(),
		), nil
	}

	nAccountRecord := types.NewAccountRecord(fact.Sender())
	nAmount := ctypes.NewAmount(accountRecord.Amount(cid.String()).Big().Sub(fact.Amount().Big()), cid)
	nAccountRecord.SetAmount(cid.String(), nAmount)
	nAccountRecord.SetLastTime(cid.String(), nowTime)

	if err := nAccountRecord.IsValid(nil); err != nil {
		return nil, base.NewBaseOperationProcessReasonError("invalid service design, %q; %w", fact.Contract(), err), nil
	}

	var sts []base.StateMergeValue // nolint:prealloc
	sts = append(sts, cstate.NewStateMergeValue(
		state.AccountRecordStateKey(fact.Contract().String(), fact.Sender().String()),
		state.NewAccountRecordStateValue(nAccountRecord),
	))

	sts = append(
		sts,
		common.NewBaseStateMergeValue(
			currency.BalanceStateKey(fact.Contract(), fact.Amount().Currency()),
			currency.NewDeductBalanceStateValue(fact.Amount()),
			func(height base.Height, st base.State) base.StateValueMerger {
				return currency.NewBalanceStateValueMerger(
					height, currency.BalanceStateKey(fact.Contract(), fact.Amount().Currency()),
					fact.Amount().Currency(), st,
				)
			}),
	)

	sts = append(sts, common.NewBaseStateMergeValue(
		currency.BalanceStateKey(fact.Receiver(), fact.Amount().Currency()),
		currency.NewAddBalanceStateValue(fact.Amount()),
		func(height base.Height, st base.State) base.StateValueMerger {
			return currency.NewBalanceStateValueMerger(height,
				currency.BalanceStateKey(fact.Receiver(), fact.Amount().Currency()),
				fact.Amount().Currency(), st,
			)
		},
	))

	return sts, nil, nil
}

func (opp *TransferProcessor) Close() error {
	opp.proposal = nil
	transferProcessorPool.Put(opp)

	return nil
}
