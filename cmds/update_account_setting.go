package cmds

import (
	"context"

	ccmds "github.com/imfact-labs/currency-model/app/cmds"
	"github.com/imfact-labs/mitum2/base"
	"github.com/imfact-labs/mitum2/util"
	"github.com/imfact-labs/payment-model/operation/payment"
	"github.com/pkg/errors"
)

type UpdateAccountInfoCommand struct {
	BaseCommand
	ccmds.OperationFlags
	Sender        ccmds.AddressFlag    `arg:"" name:"sender" help:"sender address" required:"true"`
	Contract      ccmds.AddressFlag    `arg:"" name:"contract" help:"contract address" required:"true"`
	TransferLimit ccmds.BigFlag        `arg:"" name:"transfer limit" help:"transfer limit" required:"true"`
	StartTime     uint64               `arg:"" name:"start time" help:"start time" required:"true"`
	EndTime       uint64               `arg:"" name:"end time" help:"end time" required:"true"`
	Duration      uint64               `arg:"" name:"duration" help:"duration" required:"true"`
	Currency      ccmds.CurrencyIDFlag `arg:"" name:"currency" help:"currency id" required:"true"`
	sender        base.Address
	contract      base.Address
}

func (cmd *UpdateAccountInfoCommand) Run(pctx context.Context) error { // nolint:dupl
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

	ccmds.PrettyPrint(cmd.Out, op)

	return nil
}

func (cmd *UpdateAccountInfoCommand) parseFlags() error {
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

func (cmd *UpdateAccountInfoCommand) createOperation() (base.Operation, error) { // nolint:dupl
	e := util.StringError("failed to create update account setting operation")

	fact := payment.NewUpdateAccountSettingFact([]byte(cmd.Token), cmd.sender, cmd.contract, cmd.TransferLimit.Big,
		cmd.StartTime, cmd.EndTime, cmd.Duration, cmd.Currency.CID)

	op, err := payment.NewUpdateAccountSetting(fact)
	if err != nil {
		return nil, e.Wrap(err)
	}
	err = op.Sign(cmd.Privatekey, cmd.NetworkID.NetworkID())
	if err != nil {
		return nil, e.Wrap(err)
	}

	return op, nil
}
