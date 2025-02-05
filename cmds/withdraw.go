package cmds

import (
	"context"
	currencycmds "github.com/ProtoconNet/mitum-currency/v3/cmds"
	"github.com/ProtoconNet/mitum-payment/operation/deposit"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
)

type WithdrawCommand struct {
	BaseCommand
	currencycmds.OperationFlags
	Sender          currencycmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract        currencycmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	DepositCurrency currencycmds.CurrencyIDFlag `arg:"" name:"deposit currency" help:"deposit currency id" required:"true"`
	Currency        currencycmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender          base.Address
	contract        base.Address
}

func (cmd *WithdrawCommand) Run(pctx context.Context) error { // nolint:dupl
	if _, err := cmd.prepare(pctx); err != nil {
		return err
	}

	if err := cmd.parseFlags(); err != nil {
		return err
	}

	op, err := cmd.createOperation()
	if err != nil {
		return err
	}

	currencycmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *WithdrawCommand) parseFlags() error {
	if err := cmd.OperationFlags.IsValid(nil); err != nil {
		return err
	}

	a, err := cmd.Sender.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid sender format, %q", cmd.Sender)
	} else {
		cmd.sender = a
	}

	a, err = cmd.Contract.Encode(cmd.Encoders.JSON())
	if err != nil {
		return errors.Wrapf(err, "invalid contract format, %q", cmd.Contract)
	} else {
		cmd.contract = a
	}

	return nil
}

func (cmd *WithdrawCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create withdraw operation")

	fact := deposit.NewWithdrawFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.DepositCurrency.CID, cmd.Currency.CID)

	op, err := deposit.NewWithdraw(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
